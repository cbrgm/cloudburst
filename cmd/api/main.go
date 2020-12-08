package main

import (
	"context"
	"fmt"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/ghodss/yaml"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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
			Name:  "file,f",
			Usage: "Path to the configuration file",
			Value: "cloudburst.yaml",
		},
		cli.StringFlag{
			Name:  "prometheus.url",
			Usage: "prometheus api endpoint",
			Value: "http://localhost:9090/api/v1",
		},
		cli.StringFlag{
			Name:  "addr",
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
		bytes, err := ioutil.ReadFile(c.String("file"))
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

		state, err := cloudburst.NewVolatileState(scrapeTargets)
		if err != nil {
			return fmt.Errorf("failed to initialize state provider: %w", err)
		}

		var gr run.Group
		// api
		{
			apiV1, err := NewV1(logger, state)
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
				ticker := time.NewTicker(time.Duration(15) * time.Second)
				for {
					select {
					case <-ticker.C:
						// do stuff
						err = query(promAPI, state)
						if err != nil {
							return err
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

func query(api prometheusv1.API, state *cloudburst.State) error {

	scrapeTargets, err := state.ListScrapeTargets()
	if err != nil {
		return err
	}

	for _, target := range scrapeTargets {
		value, _, err := api.Query(context.TODO(), target.Query, time.Now())
		if err != nil {
			return fmt.Errorf("failed to run query: %w", err)
		}

		vec := value.(model.Vector)

		for _, v := range vec {
			fmt.Printf("%.2f\n", v.Value)
		}
	}
	return nil
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
