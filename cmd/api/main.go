package main

import (
	"context"
	"fmt"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/cbrgm/cloudburst/state"
	"github.com/cbrgm/cloudburst/state/boltdb"
	"github.com/ghodss/yaml"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	flagAddr       = "addr"
	flagBoltPath   = "bolt.path"
	flagConfigFile = "file"
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

		var db State
		{
			var dbPath = c.String(flagBoltPath)

			bolt, dbClose, err := boltdb.NewDB(dbPath, scrapeTargets)
			if err != nil {
				return fmt.Errorf("failed to create bolt db: %s", err)
			}
			defer dbClose()

			db = state.NewStateWithProvider(bolt)
		}

		var gr run.Group
		// api
		{
			apiV1, err := NewV1(logger, db)
			if err != nil {
				return fmt.Errorf("failed to initialize api: %w", err)
			}

			r := chi.NewRouter()
			r.Use(Logger(logger))

			{
				r.Mount("/", apiV1)
			}

			s := http.Server{
				Addr:    c.String("addr"),
				Handler: r,
			}

			gr.Add(func() error {
				level.Info(logger).Log(
					"msg", "running HTTP API server",
					"addr", s.Addr,
				)
				return s.ListenAndServe()
			}, func(err error) {
				_ = s.Shutdown(context.TODO())
			})
		}

		// polling
		{
			client, err := api.NewClient(api.Config{Address: config.PrometheusURL})
			if err != nil {
				return fmt.Errorf("failed to create Prometheus client: %w", err)
			}
			promAPI := prometheusv1.NewAPI(client)

			gr.Add(func() error {
				ticker := time.NewTicker(time.Duration(5) * time.Second)
				for {
					select {
					case <-ticker.C:
						err = processScrapeTargets(promAPI, db)
						if err != nil {
							level.Info(logger).Log("msg", "prometheus processScrapeTargets job failed", "err", err)
						}
					}
				}
			}, func(err error) {

			})
		}

		{

		}

		if err := gr.Run(); err != nil {
			return errors.Errorf("error running: %w", err)
		}

		return nil

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
