package main

import (
	"fmt"
	apiclient "github.com/cbrgm/cloudburst/api/client/go"
	"github.com/zalando/skipper/eskip"
	"log"
	"net/url"
)

type Client struct {
	client  *apiclient.APIClient
	current map[string]*eskip.Route
}

// NewCloudburst creates a data client that parses a string of eskip Client and
// serves it for the routing package.
func NewCloudburst(prometheusURL string) (*Client, error) {

	apiURL, err := url.Parse(prometheusURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %w", err)
	}

	clientCfg := apiclient.NewConfiguration()
	clientCfg.Scheme = apiURL.Scheme
	clientCfg.Host = apiURL.Host

	client := apiclient.NewAPIClient(clientCfg)

	return &Client{
		client:  client,
		current: make(map[string]*eskip.Route),
	}, nil
}

func (r *Client) LoadAll() ([]*eskip.Route, error) {
	log.Println("loading Client")
	return nil, nil
}

func (*Client) LoadUpdate() ([]*eskip.Route, []string, error) {
	log.Println("updating Client")
	return nil, nil, nil
}
