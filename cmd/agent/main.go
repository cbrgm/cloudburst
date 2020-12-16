package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/cbrgm/cloudburst/cloudburst/convert"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/urfave/cli"
	"net/url"
	"os"
	"strconv"
	"time"

	apiclient "github.com/cbrgm/cloudburst/api/client/go"
)

const (
	flagName    = "name"
	flagApiAddr = "api.url"
	flagDebug   = "debug"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	app := cli.NewApp()
	app.Name = "cloudburst-agent"

	app.Action = agentAction(logger)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagName,
			Usage: "The name of the agent instance",
		},
		cli.StringFlag{
			Name:  flagApiAddr,
			Usage: "FQDN for the cloudburst-api to connect to, in format http://localhost:6660",
		},
		cli.BoolFlag{
			Name:  flagDebug,
			Usage: "debug mode",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log("msg", "failed running agent", "err", err)
		os.Exit(1)
	}
}

func agentAction(logger log.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		if c.String(flagName) == "" {
			return errors.New("no agent name provided, please set a name with --name flag")
		}

		agentName := c.String(flagName)

		if c.String(flagApiAddr) == "" {
			return errors.New("no api.url provided, please set the remote cloudburst-api url using --api.url flag")
		}

		apiURL, err := url.Parse(c.String(flagApiAddr))
		if err != nil {
			return fmt.Errorf("failed to parse API URL: %w", err)
		}

		clientCfg := apiclient.NewConfiguration()
		clientCfg.Scheme = apiURL.Scheme
		clientCfg.Host = apiURL.Host

		client := apiclient.NewAPIClient(clientCfg)

		var gr run.Group
		{
			if c.Bool(flagDebug) {
				var ticker = make(chan int)
				gr.Add(func() error {
					scan := bufio.NewScanner(os.Stdin)
					for scan.Scan() {
						s := scan.Text()
						queryValue, err := strconv.Atoi(s)
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
						case <-ticker:
							_ = poll(client, agentName, logger)
						}
					}
				}, func(err error) {

				})
			}

		}

		if !c.Bool(flagDebug) {
			ticker := time.NewTicker(time.Duration(5) * time.Second)
			gr.Add(func() error {
				for {
					select {
					case <-ticker.C:
						_ = poll(client, agentName, logger)
					}
				}
			}, func(err error) {

			})
		}

		if err := gr.Run(); err != nil {
			level.Error(logger).Log(
				"msg", "error running",
				"err", err,
			)
			os.Exit(1)
		}

		return nil
	}
}

func poll(client *apiclient.APIClient, agentName string, logger log.Logger) error {
	scrapeTargets, res, err := client.TargetsApi.ListScrapeTargets(context.TODO()).Execute()
	if err != nil || res.StatusCode != 200 {
		level.Info(logger).Log("msg", "failed retrieving scrapeTargets", "err", err)
	}

	targets := convert.APIClientToScrapeTargets(scrapeTargets)

	err = processScrapeTargets(client, agentName, targets)
	if err != nil {
		return err
	}
	return nil
}

func processScrapeTargets(client *apiclient.APIClient, agentName string, scrapeTargets []*cloudburst.ScrapeTarget) error {
	for _, target := range scrapeTargets {
		err := processScrapeTarget(client, agentName, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func processScrapeTarget(client *apiclient.APIClient, agentName string, scrapeTarget *cloudburst.ScrapeTarget) error {

	items, resp, err := client.InstancesApi.GetInstances(context.TODO(), scrapeTarget.Name).Execute()
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("failed to receive items for scrapeTarget %s", scrapeTarget.Name)
	}

	instances := convert.APIClientToInstances(items)

	terminate := cloudburst.GetInstancesByActiveStatus(instances, false)
	pending := cloudburst.GetInstancesByStatus(instances, cloudburst.Pending)

	for _, item := range pending {
		item.Status = cloudburst.InstanceStatus{
			Agent:   agentName,
			Status:  cloudburst.Progress,
			Started: time.Now(),
		}
	}

	items, resp, err = client.InstancesApi.SaveInstances(context.TODO(), scrapeTarget.Name).
		Instance(convert.InstancesToAPIClient(pending)).
		Execute()

	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("failed to update items for scrapeTarget %s", scrapeTarget.Name)
	}

	progress := convert.APIClientToInstances(items)

	// create/delete items
	for _, instance := range progress {
		instance.Status.Status = cloudburst.Running
		instance.Endpoint = "http://localhost:9998"
		time.Sleep(time.Duration(1) * time.Second)
	}

	for _, instance := range terminate {
		instance.Status.Status = cloudburst.Terminated
		time.Sleep(time.Duration(1) * time.Second)
	}

	result := append(progress, terminate...)

	// update result
	_, resp, err = client.InstancesApi.SaveInstances(context.TODO(), scrapeTarget.Name).
		Instance(convert.InstancesToAPIClient(result)).
		Execute()

	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("failed to update items for scrapeTarget %s", scrapeTarget.Name)
	}
	return nil
}
