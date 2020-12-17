package convert

import (
	openapi "github.com/cbrgm/cloudburst/api/server/go/go"
	"github.com/cbrgm/cloudburst/cloudburst"
)

func InstanceToOpenAPI(i *cloudburst.Instance) openapi.Instance {
	return openapi.Instance{
		Name:     i.Name,
		Endpoint: i.Endpoint,
		Active:   i.Active,
		Container: openapi.ContainerSpec{
			Name:  i.Container.Name,
			Image: i.Container.Image,
		},
		Status: openapi.InstanceStatus{
			Agent:   i.Status.Agent,
			Status:  string(i.Status.Status),
			Started: i.Status.Started,
		},
	}
}

func OpenAPItoInstance(i openapi.Instance) *cloudburst.Instance {
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

	return &cloudburst.Instance{
		Name:     i.Name,
		Endpoint: i.Endpoint,
		Container: cloudburst.ContainerSpec{
			Name:  i.Container.Name,
			Image: i.Container.Image,
		},
		Active: i.Active,
		Status: cloudburst.InstanceStatus{
			Agent:   i.Status.Agent,
			Status:  status,
			Started: i.Status.Started,
		},
	}
}

func InstanceEventToOpenAPI(e cloudburst.InstanceEvent) openapi.InstanceEvent {
	return openapi.InstanceEvent{
		Type:   string(e.EventType),
		Target: e.ScrapeTarget,
		Data:   InstanceToOpenAPI(e.Instance),
	}
}

func ScrapeTargetToOpenAPI(s *cloudburst.ScrapeTarget) openapi.ScrapeTarget {
	return openapi.ScrapeTarget{
		Name:        s.Name,
		Description: s.Description,
		Path:        s.Path,
		Query:       s.Query,
		InstanceSpec: openapi.InstanceSpec{
			Container: openapi.ContainerSpec{
				Name:  s.InstanceSpec.Container.Name,
				Image: s.InstanceSpec.Container.Image,
			},
		},
		StaticSpec: openapi.StaticSpec{
			Endpoints: s.StaticSpec.Endpoints,
		},
	}
}
