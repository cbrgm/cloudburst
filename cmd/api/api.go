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
	InstanceSetter
	InstanceGetter
}

func NewV1(logger log.Logger, state State) (http.Handler, error) {

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

	var body []openapi.ScrapeTarget
	for _, st := range scrapeTargets {
		body = append(body, scrapeTargetOpenAPI(st))
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
	SaveInstances(targetName string, instances []cloudburst.Instance) ([]cloudburst.Instance, error)
}

func (i *Instances) SaveInstances(ctx context.Context, targetName string, instances []openapi.Instance) (openapi.ImplResponse, error) {
	var in []cloudburst.Instance
	for _, item := range instances {
		in = append(in, instanceCloudburst(item))
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
	GetInstancesForTarget(scrapeTarget string) ([]cloudburst.Instance, error)
}

func (s *Instances) GetInstances(ctx context.Context, targetName string) (openapi.ImplResponse, error) {
	instances, err := s.getter.GetInstancesForTarget(targetName)
	if err != nil {
		return openapi.ImplResponse{
			Code: 500,
			Body: nil,
		}, err
	}

	var body []openapi.Instance
	for _, st := range instances {
		body = append(body, instanceOpenAPI(st))
	}

	return openapi.ImplResponse{
		Code: 200,
		Body: body,
	}, nil
}

func scrapeTargetOpenAPI(s cloudburst.ScrapeTarget) openapi.ScrapeTarget {
	return openapi.ScrapeTarget{
		Name:        s.Name,
		Description: s.Description,
		Query:       s.Query,
		InstanceSpec: openapi.InstanceSpec{
			Container: openapi.ContainerSpec{
				Name:  s.InstanceSpec.Container.Name,
				Image: s.InstanceSpec.Container.Image,
			},
		},
	}
}

func instanceOpenAPI(i cloudburst.Instance) openapi.Instance {
	return openapi.Instance{
		Name:     i.Name,
		Target:   i.Target,
		Endpoint: i.Endpoint,
		Active:   i.Active,
		Status: openapi.InstanceStatus{
			Agent:   i.Status.Agent,
			Status:  string(i.Status.Status),
			Started: i.Status.Started,
		},
	}
}

func instanceCloudburst(i openapi.Instance) cloudburst.Instance {
	var status cloudburst.Status
	switch i.Status.Status {
	case "unknown":
		status = cloudburst.Unknown
	case "pending":
		status = cloudburst.Pending
	case "running":
		status = cloudburst.Running
	case "failure":
		status = cloudburst.Failure
	case "progress":
		status = cloudburst.Progress
	case "terminated":
		status = cloudburst.Terminated
	}

	return cloudburst.Instance{
		Name:     i.Name,
		Endpoint: i.Endpoint,
		Target:   i.Target,
		Active:   i.Active,
		Status: cloudburst.InstanceStatus{
			Agent:   i.Status.Agent,
			Status:  status,
			Started: i.Status.Started,
		},
	}
}
