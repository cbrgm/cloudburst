package convert

import (
	apiclient "github.com/cbrgm/cloudburst/api/client/go"
	"github.com/cbrgm/cloudburst/cloudburst"
)

func APIClientToScrapeTargets(scrapeTargets []apiclient.ScrapeTarget) []*cloudburst.ScrapeTarget {
	res := []*cloudburst.ScrapeTarget{}
	for _, scrapeTarget := range scrapeTargets {
		res = append(res, APIClientToScrapeTarget(scrapeTarget))
	}
	return res
}

func APIClientToScrapeTarget(s apiclient.ScrapeTarget) *cloudburst.ScrapeTarget {
	return &cloudburst.ScrapeTarget{
		Name:        s.Name,
		Path:        s.Path,
		Description: s.Description,
		Query:       s.Query,
		ProviderSpec: cloudburst.ProviderSpec{
			Weights: s.ProviderSpec.Weights,
		},
		InstanceSpec: cloudburst.InstanceSpec{
			Container: cloudburst.ContainerSpec{
				Name:  s.InstanceSpec.Container.Name,
				Image: s.InstanceSpec.Container.Image,
			},
		},
		StaticSpec: cloudburst.StaticSpec{
			Endpoints: s.StaticSpec.Endpoints,
		},
	}
}

func InstanceToAPIClient(in *cloudburst.Instance) apiclient.Instance {
	return apiclient.Instance{
		Name:     in.Name,
		Endpoint: in.Endpoint,
		Provider: in.Provider,
		Active:   in.Active,
		Container: apiclient.ContainerSpec{
			Name:  in.Container.Name,
			Image: in.Container.Image,
		},
		Status: apiclient.InstanceStatus{
			Agent:   in.Status.Agent,
			Status:  string(in.Status.Status),
			Started: in.Status.Started,
		},
	}
}

func InstancesToAPIClient(instances []*cloudburst.Instance) []apiclient.Instance {
	res := []apiclient.Instance{}
	for _, instance := range instances {
		res = append(res, InstanceToAPIClient(instance))
	}
	return res
}

func APIClientToInstance(in apiclient.Instance) *cloudburst.Instance {
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

func APIClientToInstances(instances []apiclient.Instance) []*cloudburst.Instance {
	res := []*cloudburst.Instance{}
	for _, instance := range instances {
		res = append(res, APIClientToInstance(instance))
	}
	return res
}
