package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"os"

	"github.com/zalando/skipper"
	"github.com/zalando/skipper/routing"
)

const (
	flagApiAddr = "api.url"
	flagAddr    = "addr"
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
			Name:  flagAddr,
			Usage: "The address for the public http server",
			Value: ":6661",
		},
		cli.StringFlag{
			Name:  flagApiAddr,
			Usage: "FQDN for the cloudburst-api to connect to, in format http://localhost:6660",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log("msg", "failed running proxy", "err", err)
		os.Exit(1)
	}
}

func proxyAction(logger log.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {

		if c.String(flagApiAddr) == "" {
			return errors.New("no proxy addr provided, please set an address to bind on using --addr flag")
		}

		if c.String(flagApiAddr) == "" {
			return errors.New("no api.url provided, please set the remote cloudburst-api url using --api.url flag")
		}

		client, err := NewCloudburst(c.String(flagApiAddr))
		if err != nil {
			logger.Log("msg", "failed to create client for proxy", "err", err)
		}

		var gr run.Group
		{
			gr.Add(func() error {
				level.Info(logger).Log(
					"msg", "running HTTP proxy server",
					"addr", c.String(flagAddr),
				)
				err := skipper.Run(skipper.Options{
					Address:           c.String(flagAddr),
					CustomDataClients: []routing.DataClient{client},
				})
				return err
			}, func(err error) {

			})
		}

		if err := gr.Run(); err != nil {
			return errors.Errorf("error running: %w", err)
		}
		return nil
	}
}
