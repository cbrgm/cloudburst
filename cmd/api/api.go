package main

import (
	"context"
	openapi "github.com/cbrgm/cloudburst/api/server/go/go"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/cbrgm/cloudburst/cloudburst/convert"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"net/http"
)

type State interface {
	ScrapeTargetLister
	InstanceSetter
	InstanceGetter
}

func NewV1(logger log.Logger, state State, events Events) (http.Handler, error) {

	routes := []openapi.Router{
		openapi.NewTargetsApiController(&ScrapeTargets{
			lister: state,
		}),
		openapi.NewInstancesApiController(&Instances{
			setter: state,
			getter: state,
		}),
	}

	router := mux.NewRouter().StrictSlash(true)

	// TODO: add prometheus instrumentation
	for _, api := range routes {
		for _, route := range api.Routes() {
			router.Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(route.HandlerFunc)
		}
	}

	router.Methods(http.MethodGet).
		Path("/api/v1/instances/events").
		HandlerFunc(instanceEventsHandler(logger, events))

	return router, nil
}

type ScrapeTargets struct {
	lister ScrapeTargetLister
}

type ScrapeTargetLister interface {
	ListScrapeTargets() ([]*cloudburst.ScrapeTarget, error)
}

func (s *ScrapeTargets) ListScrapeTargets(ctx context.Context) (openapi.ImplResponse, error) {
	scrapeTargets, err := s.lister.ListScrapeTargets()
	if err != nil {
		return openapi.ImplResponse{
			Code: 500,
			Body: nil,
		}, err
	}

	var body []openapi.ScrapeTarget
	for _, st := range scrapeTargets {
		body = append(body, convert.ScrapeTargetToOpenAPI(st))
	}

	return openapi.ImplResponse{
		Code: 200,
		Body: body,
	}, nil
}

type Instances struct {
	setter InstanceSetter
	getter InstanceGetter
}

type InstanceSetter interface {
	SaveInstances(targetName string, instances []*cloudburst.Instance) ([]*cloudburst.Instance, error)
	SaveInstance(targetName string, instances *cloudburst.Instance) (*cloudburst.Instance, error)
}

func (i *Instances) SaveInstances(ctx context.Context, targetName string, instances []openapi.Instance) (openapi.ImplResponse, error) {
	var in []*cloudburst.Instance
	for _, item := range instances {
		in = append(in, convert.OpenAPItoInstance(item))
	}

	body, err := i.setter.SaveInstances(targetName, in)
	if err != nil {
		return openapi.ImplResponse{
			Code: 500,
			Body: nil,
		}, err
	}

	return openapi.ImplResponse{
		Code: 200,
		Body: body,
	}, nil
}

type InstanceGetter interface {
	GetInstances(scrapeTarget string) ([]*cloudburst.Instance, error)
}

func (s *Instances) GetInstances(ctx context.Context, targetName string) (openapi.ImplResponse, error) {
	instances, err := s.getter.GetInstances(targetName)
	if err != nil {
		return openapi.ImplResponse{
			Code: 500,
			Body: nil,
		}, err
	}

	var body []openapi.Instance
	for _, st := range instances {
		body = append(body, convert.InstanceToOpenAPI(st))
	}

	return openapi.ImplResponse{
		Code: 200,
		Body: body,
	}, nil
}
