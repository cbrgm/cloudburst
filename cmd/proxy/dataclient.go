package main

import (
	"context"
	"fmt"
	apiclient "github.com/cbrgm/cloudburst/api/client/go"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/cbrgm/cloudburst/cloudburst/convert"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/zalando/skipper/eskip"
	"net/http"
	"net/url"
)

type Client struct {
	client *apiclient.APIClient
	old    map[string]*eskip.Route
	logger log.Logger
}

// NewCloudburst creates a data client that parses a string of eskip Client and
// serves it for the routing package.
func NewCloudburst(logger log.Logger, prometheusURL string) (*Client, error) {
	apiURL, err := url.Parse(prometheusURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %v", err)
	}

	clientCfg := apiclient.NewConfiguration()
	clientCfg.Scheme = apiURL.Scheme
	clientCfg.Host = apiURL.Host

	client := apiclient.NewAPIClient(clientCfg)

	return &Client{
		client: client,
		logger: logger,
		old:    make(map[string]*eskip.Route),
	}, nil
}

func (c *Client) loadAndConvert() ([]*eskip.Route, error) {
	res, resp, err := c.client.TargetsApi.ListScrapeTargets(context.TODO()).Execute()
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	scrapeTargets := convert.APIClientToScrapeTargets(res)

	r, err := c.routesForScrapeTargets(scrapeTargets)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) routesForScrapeTargets(scrapeTargets []*cloudburst.ScrapeTarget) ([]*eskip.Route, error) {
	var res []*eskip.Route
	for _, target := range scrapeTargets {
		route, err := c.routeForScrapeTarget(target)
		if err != nil {
			return res, err
		}
		res = append(res, route)
	}
	return res, nil
}

func (c *Client) routeForScrapeTarget(scrapeTarget *cloudburst.ScrapeTarget) (*eskip.Route, error) {
	res, resp, err := c.client.InstancesApi.GetInstances(context.TODO(), scrapeTarget.Name).Execute()
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	instances := convert.APIClientToInstances(res)

	lbEndpoints := []string{}
	for _, instance := range instances {
		if instance.Endpoint != "" {
			lbEndpoints = append(lbEndpoints, instance.Endpoint)
		}
	}
	lbEndpoints = append(lbEndpoints, scrapeTarget.StaticSpec.Endpoints...)

	return &eskip.Route{
		Id: scrapeTarget.Name,
		Predicates: []*eskip.Predicate{
			{
				Name: "Path",
				Args: []interface{}{
					scrapeTarget.Path,
				},
			},
		},
		Filters:     []*eskip.Filter{},
		BackendType: 4,
		Backend:     "",
		LBAlgorithm: "random",
		LBEndpoints: lbEndpoints,
	}, nil
}

func mapRoutes(r []*eskip.Route) map[string]*eskip.Route {
	m := make(map[string]*eskip.Route)
	for _, ri := range r {
		m[ri.Id] = ri
	}

	return m
}

func (c *Client) LoadAll() ([]*eskip.Route, error) {
	level.Debug(c.logger).Log("msg", "loading routes")

	r, err := c.loadAndConvert()
	if err != nil {
		return nil, fmt.Errorf("failed to load routes from cloudburst-api: %v", err)
	}

	c.old = mapRoutes(r)
	level.Debug(c.logger).Log("msg", "all routes loaded and mapped")

	return r, nil
}

func (c *Client) LoadUpdate() ([]*eskip.Route, []string, error) {
	level.Debug(c.logger).Log("msg", "updating routes")

	r, err := c.loadAndConvert()
	if err != nil {
		level.Error(c.logger).Log("msg", "polling for updates failed: %v", err)
		return nil, nil, fmt.Errorf("polling for updates from cloudburst-api failed: %s", err)
	}

	new := mapRoutes(r)
	level.Debug(c.logger).Log("msg", "new version of routes loaded and mapped")

	var (
		updatedRoutes []*eskip.Route
		deletedIDs    []string
	)

	for key := range c.old {
		if value, ok := new[key]; ok && value.String() != c.old[key].String() {
			updatedRoutes = append(updatedRoutes, value)
		} else if !ok {
			deletedIDs = append(deletedIDs, key)
		}
	}

	for id, r := range new {
		if _, ok := c.old[id]; !ok {
			updatedRoutes = append(updatedRoutes, r)
		}
	}

	if len(updatedRoutes) > 0 || len(deletedIDs) > 0 {
		level.Info(c.logger).Log("msg", "diff taken", "updates",  len(updatedRoutes), "deletions", len(deletedIDs))
	}

	c.old = new
	return updatedRoutes, deletedIDs, nil
}
