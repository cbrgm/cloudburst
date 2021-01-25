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

	instancesGauge *prometheus.GaugeVec
	demandGauge *prometheus.GaugeVec
	queryGauge *prometheus.GaugeVec
}

func NewAutoScaler(r *prometheus.Registry, state State) Autoscaler {

	metrics := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cloudburst_api_instances_total",
		Help: "instances total",
	}, []string{"target", "status"})

	demand := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cloudburst_api_demand",
		Help: "autoscaler demand",
	}, []string{"target"})

	query := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cloudburst_api_query",
		Help: "autoscaler query",
	}, []string{"target"})

	r.MustRegister(metrics)
	r.MustRegister(demand)
	r.MustRegister(query)

	threshold := newThreshold(0, 0)
	return &autoscaler{
		state:          state,
		requester:      newRequester(state, threshold),
		instancesGauge: metrics,
		demandGauge: demand,
		queryGauge: query,
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


	metrics := make(map[Status]float64)
	for _, instance := range instances {
		if val, ok := metrics[instance.Status.Status]; ok {
			metrics[instance.Status.Status] = val + 1
		} else {
			metrics[instance.Status.Status] = 1
		}
	}
	s.instancesGauge.WithLabelValues(scrapeTarget.Name, string(Pending)).Set(metrics[Pending])
	s.instancesGauge.WithLabelValues(scrapeTarget.Name, string(Progress)).Set(metrics[Progress])
	s.instancesGauge.WithLabelValues(scrapeTarget.Name, string(Running)).Set(metrics[Running])
	s.instancesGauge.WithLabelValues(scrapeTarget.Name, string(Terminated)).Set(metrics[Terminated])
	s.instancesGauge.WithLabelValues(scrapeTarget.Name, string(Failure)).Set(metrics[Failure])

	demand := s.calculateDemand(scrapeTarget, instances, queryResult)

	s.demandGauge.WithLabelValues(scrapeTarget.Name).Set(float64(demand.Result))
	s.queryGauge.WithLabelValues(scrapeTarget.Name).Set(queryResult)

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

	sumInternal := float64(len(scrapeTarget.StaticSpec.Endpoints))
	sumExternal := float64(CountInstancesByStatus(instances, Running))
	sumTerminating := float64(CountActiveInstances(instances, false))
	sumProgressActive := float64(CountActiveInstances(GetInstancesByStatus(instances, Progress), true))

	// queryResult == (sum(rate(example_sorting_requests_total[15s])) / 15) => CONFIG
	sumEffective := math.Round((queryResult -1) * (sumInternal + sumExternal) + 0.5)
	demand := (sumEffective + sumTerminating) - sumProgressActive

	return instanceDemand{
		Result: int(demand),
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
