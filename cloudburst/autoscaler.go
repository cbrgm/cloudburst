package cloudburst

import (
	"github.com/prometheus/client_golang/prometheus"
	"math"
)

type Autoscaler interface {
	Scale(scrapeTarget *ScrapeTarget, queryResult float64) error
}

type autoscaler struct {
	state     State
	requester *requester

	instancesCounter *prometheus.CounterVec
}

func NewAutoScaler(r *prometheus.Registry, state State) Autoscaler {

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "instances_total",
		Help: "instances total",
	}, []string{"target", "status"})

	r.MustRegister(counter)

	threshold := newThreshold(0, -1)
	return &autoscaler{
		state:            state,
		requester:        newRequester(state, threshold),
		instancesCounter: counter,
	}
}

func (s *autoscaler) Scale(scrapeTarget *ScrapeTarget, queryResult float64) error {

	// TODO: don't call state twice
	err := s.cleanupTerminatedInstances(scrapeTarget)
	if err != nil {
		return err
	}

	instances, err := s.state.GetInstances(scrapeTarget.Name)
	if err != nil {
		return err
	}

	for _, instance := range instances {
		s.instancesCounter.WithLabelValues(scrapeTarget.Name, string(instance.Status.Status)).Inc()
	}

	demand := s.calculateDemand(scrapeTarget, instances, queryResult)
	err = s.processDemand(demand, instances, scrapeTarget)
	if err != nil {
		return err
	}

	return nil
}

type instanceDemand struct {
	Result int
}

// CalculateDemand calculates the queryResult for Instance objects. The provided instances slice is a list of all instances
// to be filtered for calculation. The provided queryResult value is the result of a metric query.
// CalculateDemand returns instanceDemand.
func (s *autoscaler) calculateDemand(scrapeTarget *ScrapeTarget, instances []*Instance, queryResult float64) instanceDemand {

	sumRunningInstancesExternAndIntern := float64(CountInstancesByStatus(instances, Running) + len(scrapeTarget.StaticSpec.Endpoints))
	queryResult = math.Round(((queryResult - 1) * sumRunningInstancesExternAndIntern) + 0.5)

	sumTerminating := CountActiveInstances(instances, false)
	progress := GetInstancesByStatus(instances, Progress)
	sumProgressStart := CountActiveInstances(progress, true)

	var demand = (int(queryResult) + sumTerminating) - sumProgressStart

	return instanceDemand{
		Result: demand,
	}
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

func (s *autoscaler) processDemand(demand instanceDemand, instances []*Instance, scrapeTarget *ScrapeTarget) error {
	return s.requester.ProcessDemand(demand, instances, scrapeTarget)
}
