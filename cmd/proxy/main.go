package main

import (
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/urfave/cli"
	"os"

	"github.com/zalando/skipper"
	"github.com/zalando/skipper/routing"
)

const (
	flagApiAddr = "api.url"
)
func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	app := cli.NewApp()
	app.Name = "cloudburst-proxy"

	app.Action = proxyAction(logger)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagApiAddr,
			Usage: "FQDN for the cloudburst-api to connect to, in format http://localhost:6660",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log("msg", "failed running agent", "err", err)
		os.Exit(1)
	}
}

func proxyAction(logger log.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {

		if c.String(flagApiAddr) == "" {
			return errors.New("no api.url provided, please set the remote cloudburst-api url using --api.url flag")
		}

		client, err := NewCloudburst(c.String(flagApiAddr))
		if err != nil {
			logger.Log("msg", "failed to create client for proxy", "err", err)
		}

		skipper.Run(skipper.Options{
			Address:           ":9090",
			CustomDataClients: []routing.DataClient{client},
		})

		return nil
	}
}
