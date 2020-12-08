package main

import (
	"context"
	openapi "github.com/cbrgm/cloudburst/api/server/go/go"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"net/http"
)

type State interface {
	ScrapeTargetLister
	InstanceUpdater
}

func NewV1(logger log.Logger, state State) (http.Handler, error) {

	routes := []openapi.Router{
		openapi.NewTargetsApiController(&ScrapeTargets{
			lister: state,
		}),
		openapi.NewInstancesApiController(&Instances{
			updater: state,
		}),
	}

	router := mux.NewRouter().StrictSlash(true)

	// TODO: add prometheus instrumentation
	for _, api := range routes {
		for _, route := range api.Routes() {
			router.Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name)
		}
	}

	return router, nil
}

type ScrapeTargets struct {
	lister ScrapeTargetLister
}

type ScrapeTargetLister interface {
	ListScrapeTargets() ([]cloudburst.ScrapeTarget, error)
}

func (s *ScrapeTargets) ListScrapeTargets(ctx context.Context) (openapi.ImplResponse, error) {
	scrapeTargets, err := s.lister.ListScrapeTargets()
	if err != nil {
		return openapi.ImplResponse{
			Code: 500,
			Body: nil,
		}, err
	}

	var targets []openapi.ScrapeTarget
	for _, st := range scrapeTargets {
		targets = append(targets, scrapeTargetOpenAPI(st))
	}

	return openapi.ImplResponse{
		Code: 200,
		Body: targets,
	}, nil
}

type Instances struct {
	updater InstanceUpdater
}

type InstanceUpdater interface {
	UpdateInstances(targetName string, instances []cloudburst.Instance) ([]cloudburst.Instance, error)
}

func (i *Instances) UpdateInstances(ctx context.Context, targetName string, instances []openapi.Instance) (openapi.ImplResponse, error) {
	var in []cloudburst.Instance
	for _, item := range instances {
		in = append(in, instanceCloudburst(item))
	}

	res, err := i.updater.UpdateInstances(targetName, in)
	if err != nil {
		return openapi.ImplResponse{
			Code: 500,
			Body: nil,
		}, err
	}

	return openapi.ImplResponse{
		Code: 200,
		Body: res,
	}, nil
}

func scrapeTargetOpenAPI(s cloudburst.ScrapeTarget) openapi.ScrapeTarget {

	var instances []openapi.Instance
	for _, item := range s.Instances {
		instances = append(instances, instanceOpenAPI(item))
	}

	return openapi.ScrapeTarget{
		Name:        s.Name,
		Description: s.Description,
		Query:       s.Query,
		Value:       s.Value,
		InstanceSpec: openapi.InstanceSpec{
			Container: openapi.ContainerSpec{
				Name:  s.InstanceSpec.Container.Name,
				Image: s.InstanceSpec.Container.Image,
			},
		},
		Instances: instances,
	}
}

func instanceOpenAPI(i cloudburst.Instance) openapi.Instance {
	return openapi.Instance{
		Name:     i.Name,
		Endpoint: i.Endpoint,
		Active:   i.Active,
		Status: openapi.InstanceStatus{
			Agent:   i.Status.Agent,
			Status:  i.Status.Status,
			Started: i.Status.Started,
		},
	}
}

func instanceCloudburst(i openapi.Instance) cloudburst.Instance {
	return cloudburst.Instance{
		Name:     i.Name,
		Endpoint: i.Endpoint,
		Active:   i.Active,
		Status: cloudburst.InstanceStatus{
			Agent:   i.Status.Agent,
			Status:  i.Status.Status,
			Started: i.Status.Started,
		},
	}
}
