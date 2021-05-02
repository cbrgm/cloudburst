package main

import (
	"context"
	"fmt"
	openapi "github.com/cbrgm/cloudburst/api/server/go/go"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/cbrgm/cloudburst/cloudburst/convert"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type State interface {
	ScrapeTargetLister
	InstanceSetter
	InstanceGetter
}

type ApiMetrics interface {
	MeasureOpenApiRequestDuration(code, method, name string, start time.Time)
	IncOpenApiRequestsTotal(code, method, name string)
	MeasureApiEventDuration(error string, start time.Time)
	IncEventSubscribers()
	DecEventSubscribers()
}

func NewV1(logger log.Logger, apiMetrics ApiMetrics, state State, events Events) (http.Handler, error) {
	instrument := instrument(apiMetrics)
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

	for _, api := range routes {
		for _, route := range api.Routes() {
			router.Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(instrument(
					route.HandlerFunc,
					route.Name,
				))
		}
	}

	router.Methods(http.MethodGet).
		Path("/api/v1/instances/events").
		HandlerFunc(instanceEventsHandler(logger, apiMetrics, events))

	return router, nil
}

func instrument(m ApiMetrics) func(next http.Handler, name string) http.Handler {
	return func(next http.Handler, name string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			m.MeasureOpenApiRequestDuration(
				fmt.Sprintf("%d", ww.Status()),
				r.Method,
				name,
				start,
			)
		})
	}
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

	body := []openapi.ScrapeTarget{}
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
	in := []*cloudburst.Instance{}
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

	body := []openapi.Instance{}
	for _, st := range instances {
		body = append(body, convert.InstanceToOpenAPI(st))
	}

	return openapi.ImplResponse{
		Code: 200,
		Body: body,
	}, nil
}
