package cloudburst

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Autoscaler interface {
	Scale(scrapeTarget *ScrapeTarget, metricValue float64) error
}

type autoscaler struct {
	state     State
	requester *requester
	scale     ScalingFunc
}

func NewInstrumentedAutoScaler(r *prometheus.Registry, scalingFunc ScalingFunc, state State) Autoscaler {
	threshold := newThreshold(0, 0)
	return &autoscaler{
		state:     state,
		requester: newRequester(state, threshold),
		scale:     InstrumentedScalingFunc(r, scalingFunc),
	}
}

func NewAutoScaler(scale ScalingFunc, state State) Autoscaler {
	threshold := newThreshold(0, 0)
	return &autoscaler{
		state:     state,
		requester: newRequester(state, threshold),
		scale:     scale,
	}
}

func (s *autoscaler) Scale(scrapeTarget *ScrapeTarget, metricValue float64) error {
	// TODO: don't call state twice

	err := s.cleanupTerminatedInstances(scrapeTarget)
	if err != nil {
		return err
	}

	instances, err := s.state.GetInstances(scrapeTarget.Name)
	if err != nil {
		return err
	}

	demand := s.scale(scrapeTarget, instances, metricValue)

	err = s.processDemand(demand, scrapeTarget)
	if err != nil {
		return err
	}

	return nil
}

func (s *autoscaler) cleanupTerminatedInstances(target *ScrapeTarget) error {
	instances, err := s.state.GetInstances(target.Name)
	if err != nil {
		return err
	}

	for _, instance := range instances {
		if instance.Status.Status == Terminated || instance.Status.Status == Failure {
			err := s.state.RemoveInstance(target.Name, instance)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *autoscaler) processDemand(demand ScalingResult, scrapeTarget *ScrapeTarget) error {
	return s.requester.ProcessDemand(demand, scrapeTarget)
}
