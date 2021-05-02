package convert

import (
	openapi "github.com/cbrgm/cloudburst/api/server/go/go"
	"github.com/cbrgm/cloudburst/cloudburst"
)

func InstanceToOpenAPI(in *cloudburst.Instance) openapi.Instance {
	return openapi.Instance{
		Name:     in.Name,
		Endpoint: in.Endpoint,
		Provider: in.Provider,
		Active:   in.Active,
		Container: openapi.ContainerSpec{
			Name:  in.Container.Name,
			Image: in.Container.Image,
		},
		Status: openapi.InstanceStatus{
			Agent:   in.Status.Agent,
			Status:  string(in.Status.Status),
			Started: in.Status.Started,
		},
	}
}

func InstancesToOpenAPI(instances []*cloudburst.Instance) []openapi.Instance {
	res := []openapi.Instance{}
	for _, instance := range instances {
		res = append(res, InstanceToOpenAPI(instance))
	}
	return res
}

func OpenAPItoInstance(in openapi.Instance) *cloudburst.Instance {
	var status cloudburst.Status
	switch in.Status.Status {
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
		Name:     in.Name,
		Endpoint: in.Endpoint,
		Provider: in.Provider,
		Active:   in.Active,
		Container: cloudburst.ContainerSpec{
			Name:  in.Container.Name,
			Image: in.Container.Image,
		},
		Status: cloudburst.InstanceStatus{
			Agent:   in.Status.Agent,
			Status:  status,
			Started: in.Status.Started,
		},
	}
}

func InstanceEventToOpenAPI(e cloudburst.InstanceEvent) openapi.InstanceEvent {
	return openapi.InstanceEvent{
		Type:   string(e.EventType),
		Target: e.ScrapeTarget,
		Data:   InstancesToOpenAPI(e.Instances),
	}
}

func ScrapeTargetToOpenAPI(s *cloudburst.ScrapeTarget) openapi.ScrapeTarget {
	return openapi.ScrapeTarget{
		Name:        s.Name,
		Description: s.Description,
		Path:        s.Path,
		Query:       s.Query,
		ProviderSpec: openapi.ProviderSpec{
			Weights: s.ProviderSpec.Weights,
		},
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
