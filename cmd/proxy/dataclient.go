package main

import (
	apiclient "github.com/cbrgm/cloudburst/api/client/go"
	"github.com/zalando/skipper/eskip"
	"github.com/zalando/skipper/routing"
	"log"
)

type routes struct {
	parsed []*eskip.Route
}

// NewCloudburst creates a data client that parses a string of eskip routes and
// serves it for the routing package.
func NewCloudburst() (routing.DataClient, error) {
	route := `* -> inlineContent("Hello, world!") -> <shunt>`
	parsed, err := eskip.Parse(route)
	if err != nil {
		return nil, err
	}

	return &routes{parsed: parsed}, nil
}

func (r *routes) LoadAll() ([]*eskip.Route, error) {
	log.Println("loading routes")
	return r.parsed, nil
}

func (*routes) LoadUpdate() ([]*eskip.Route, []string, error) {
	log.Println("updating routes")
	return nil, nil, nil
}
