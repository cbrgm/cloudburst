package api

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"time"
)

type configuration struct {
	PrometheusURL string `json:"prometheus_url"`
	ScrapeTargets []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Query       string `json:"query"`
		Spec        struct {
			Container struct {
				Name  string `json:"name"`
				Image string `json:"image"`
			} `json:"container"`
		} `json:"spec"`
	} `json:"targets"`
}

type ScrapeTarget struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Query       string       `json:"query"`
	Value       float64      `json:"value"`
	Spec        InstanceSpec `json:"spec"`
	Instances   []Instance   `json:"instances"`
}

type Instance struct {
	Name     string
	Provider string
	Endpoint string
	Status   string
}

type InstanceSpec struct {
	Container ContainerSpec `json:"container"`
}

type ContainerSpec struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

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

		var config configuration
		if err := yaml.Unmarshal(bytes, &config); err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		services, err := getTargetsFromConfig(config)
		if err != nil {
			return fmt.Errorf("failed to parse services from config: %w", err)
		}

		var gr run.Group
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
						// do stuff
						err = query(promAPI, services)
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

func query(api prometheusv1.API, targets []ScrapeTarget) error {
	return nil
}

func getTargetsFromConfig(config configuration) ([]ScrapeTarget, error) {
	var res []ScrapeTarget
	for _, svc := range config.Targets {
		res = append(res, ScrapeTarget{
			Name:        svc.Name,
			Description: svc.Description,
			Query:       svc.Query,
			Value:       0,
			Spec: InstanceSpec{
				Container: ContainerSpec{
					Name:  svc.Spec.Container.Name,
					Image: svc.Spec.Container.Image,
				},
			},
			Instances: []Instance{},
		})
	}
	return res, nil
}
