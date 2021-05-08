package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/cbrgm/cloudburst/autoscaler"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/cbrgm/cloudburst/metrics"
	"github.com/cbrgm/cloudburst/state/boltdb"
	"github.com/ghodss/yaml"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const (
	flagAddr            = "addr"
	flagInternalAddr    = "internal.addr"
	flagBoltPath        = "bolt.path"
	flagDebug           = "debug"
	flagConfigFile      = "file"
	flagUIAssets        = "ui.assets"
	flagPollingInterval = "polling.interval"
)

func main() {

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	app := cli.NewApp()
	app.Name = "cloudburst-api"

	app.Action = apiAction(logger)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagConfigFile,
			Usage: "Path to the configuration file",
			Value: "cloudburst.yaml",
		},
		cli.StringFlag{
			Name:  flagBoltPath,
			Value: "./development/data",
		},
		cli.StringFlag{
			Name:  flagAddr,
			Usage: "The address for the public http server",
			Value: ":6660",
		},
		cli.StringFlag{
			Name:  flagInternalAddr,
			Usage: "The internal address for the public http server",
			Value: ":6661",
		},
		cli.StringFlag{
			Name:  flagUIAssets,
			Usage: "The path to the ui assets",
			Value: "./ui",
		},
		cli.IntFlag{
			Name:  flagPollingInterval,
			Usage: "The interval in seconds to check slo rules",
			Value: 10,
		},
		cli.BoolFlag{
			Name:  flagDebug,
			Usage: "debug mode",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log("msg", "failed running api", "err", err)
		os.Exit(1)
	}
}

func apiAction(logger log.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		bytes, err := ioutil.ReadFile(c.String(flagConfigFile))
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		var config cloudburst.Configuration
		if err := yaml.Unmarshal(bytes, &config); err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		scrapeTargets, err := cloudburst.ParseConfiguration(config)
		if err != nil {
			return fmt.Errorf("failed to parse scrapeTargets from config: %w", err)
		}

		events := cloudburst.NewEvents()

		var state autoscaler.State
		{
			var dbPath = c.String(flagBoltPath)
			bolt, dbClose, err := boltdb.NewDB(dbPath, scrapeTargets)
			if err != nil {
				return fmt.Errorf("failed to create bolt state: %s", err)
			}
			defer dbClose()

			boltEvents := boltdb.NewEvents(bolt, events)
			state = boltEvents
		}

		metricOptions := metrics.Options{
			Enabled:              true,
			Prefix:               "",
			EnableProfile:        true,
			EnableRuntimeMetrics: true,
		}
		scalingMetrics := metrics.NewDefaultPrometheus()
		ctx, cancel := context.WithCancel(context.Background())

		var gr run.Group
		// api
		{
			apiV1, err := NewV1(logger, scalingMetrics, state, events)
			if err != nil {
				return fmt.Errorf("failed to initialize api: %w", err)
			}

			r := chi.NewRouter()
			r.Use(Logger(logger))
			{
				r.Mount("/", apiV1)
			}
			{
				directory := c.String(flagUIAssets)
				if _, err := os.Stat(directory); os.IsNotExist(err) {
					return fmt.Errorf("assets path not found: %s", directory)
				}

				// serve ui
				r.Get("/", file(directory, "index.html"))
				r.Get("/bulma.min.css", file(directory, "bulma.min.css"))
				r.Get("/bundle.js", file(directory, "bundle.js"))
				r.Get("/bundle.js.map", file(directory, "bundle.js.map"))
				r.NotFound(file(directory, "index.html"))
			}

			s := http.Server{
				Addr:    c.String("addr"),
				Handler: r,
			}

			gr.Add(func() error {
				level.Info(logger).Log(
					"msg", "running public HTTP API server",
					"addr", s.Addr,
				)
				return s.ListenAndServe()
			}, func(err error) {
				_ = s.Shutdown(context.TODO())
			})
		}
		{
			r := chi.NewRouter()
			r.Mount("/metrics", metrics.HandlerFor(scalingMetrics, metricOptions))

			s := http.Server{
				Addr:    c.String(flagInternalAddr),
				Handler: r,
			}
			gr.Add(func() error {
				level.Info(logger).Log(
					"msg", "running internal HTTP API server",
					"addr", s.Addr,
				)

				return s.ListenAndServe()
			}, func(err error) {
				_ = s.Shutdown(context.TODO())
			})
		}

		// polling
		{
			if c.Bool(flagDebug) {
				sLogger := log.With(logger, "component", "autoscaler")
				scaling, err := autoscaler.NewScaling(state, config.PrometheusURL, autoscaler.WithMetrics(scalingMetrics))
				if err != nil {
					level.Error(sLogger).Log("msg", "failed to initialize telegram bot", "err", err)
					os.Exit(2)
				}
				var ticker = make(chan float64)
				gr.Add(func() error {
					scan := bufio.NewScanner(os.Stdin)
					for scan.Scan() {
						s := scan.Text()
						queryValue, err := strconv.ParseFloat(s, 64)
						if err != nil {
							continue
						}
						ticker <- queryValue
					}
					return nil
				}, func(err error) {
					close(ticker)
				})
				gr.Add(func() error {
					for {
						select {
						case i := <-ticker:
							err = scaling.ProcessScrapeTargetWithValue(scrapeTargets[0], i)
							if err != nil {
								level.Info(logger).Log("msg", "prometheus processScrapeTargets job failed", "err", err)
							}
						}
					}
				}, func(err error) {
					cancel()
				})
			}
			if !c.Bool(flagDebug) {
				sLogger := log.With(logger, "component", "autoscaler")
				interval := time.Duration(c.Int(flagPollingInterval)) * time.Second
				scaling, err := autoscaler.NewScaling(state, config.PrometheusURL,
					autoscaler.WithLogger(sLogger),
					autoscaler.WithInterval(interval),
					autoscaler.WithMetrics(scalingMetrics),
				)
				if err != nil {
					level.Error(sLogger).Log("msg", "failed to initialize telegram bot", "err", err)
					os.Exit(2)
				}

				gr.Add(func() error {
					level.Info(sLogger).Log(
						"msg", "starting autoscaler",
					)
					return scaling.Run(ctx)
				}, func(err error) {
					cancel()
				})
			}
		}
		{
			sig := make(chan os.Signal)
			signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

			gr.Add(func() error {
				<-sig
				return nil
			}, func(err error) {
				cancel()
				close(sig)
			})
		}

		if err := gr.Run(); err != nil {
			return errors.Errorf("error running: %w", err)
		}
		return nil
	}
}

func file(directory, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(directory, name))
	}
}

// Logger returns a middleware to log HTTP requests
func Logger(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			level.Debug(logger).Log(
				"proto", r.Proto,
				"method", r.Method,
				"status", ww.Status(),
				"content", r.Header.Get("Content-Type"),
				"path", r.URL.Path,
				"duration", time.Since(start),
				"bytes", ww.BytesWritten(),
			)
		})
	}
}
