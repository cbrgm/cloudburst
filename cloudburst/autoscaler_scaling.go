package cloudburst

import (
	"github.com/prometheus/client_golang/prometheus"
	"math"
)

// ScalingFunc calculates the demand for Instance objects.
//
// instances slice is a list of all instances  to be filtered for calculation.
// metricValue value is the result of a metric query.
// CalculateDemand returns instanceDemand.
type ScalingFunc func(scrapeTarget *ScrapeTarget, instances []*Instance, metricValue float64) instanceDemand

// instanceDemand represents the calculated demand for instances
type instanceDemand struct {
	Result int // the rounded demand for new instances
}

// NewDefaultScalingFunc returns the default scaling func
func NewDefaultScalingFunc() ScalingFunc {
	return func(scrapeTarget *ScrapeTarget, instances []*Instance, metricValue float64) instanceDemand {

		sumInternal := float64(len(scrapeTarget.StaticSpec.Endpoints))
		sumExternal := float64(CountInstancesByStatus(instances, Running))
		sumTerminating := float64(CountActiveInstances(instances, false))
		sumProgressActive := float64(CountActiveInstances(GetInstancesByStatus(instances, Progress), true))

		sumEffective := math.Round((metricValue-1)*(sumInternal+sumExternal) + 0.5)
		demand := (sumEffective + sumTerminating) - sumProgressActive

		return instanceDemand{
			Result: int(demand),
		}
	}
}

func InstrumentedScalingFunc(r *prometheus.Registry, scalingFunc ScalingFunc) ScalingFunc {

	instancesGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cloudburst_api_instances_total",
		Help: "instances total",
	}, []string{"target", "status"})

	instanceDemandGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cloudburst_api_demand",
		Help: "calculated instance demand for the current scaling iteration",
	}, []string{"target"})

	metricValueGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cloudburst_api_query",
		Help: "metric value received from prometheus for the current scaling iteration",
	}, []string{"target"})

	r.MustRegister(instancesGauge)
	r.MustRegister(instanceDemandGauge)
	r.MustRegister(metricValueGauge)

	return func(scrapeTarget *ScrapeTarget, instances []*Instance, metricValue float64) instanceDemand {
		metrics := make(map[Status]float64)
		for _, instance := range instances {
			if val, ok := metrics[instance.Status.Status]; ok {
				metrics[instance.Status.Status] = val + 1
			} else {
				metrics[instance.Status.Status] = 1
			}
		}

		instancesGauge.WithLabelValues(scrapeTarget.Name, string(Pending)).Set(metrics[Pending])
		instancesGauge.WithLabelValues(scrapeTarget.Name, string(Progress)).Set(metrics[Progress])
		instancesGauge.WithLabelValues(scrapeTarget.Name, string(Running)).Set(metrics[Running])
		instancesGauge.WithLabelValues(scrapeTarget.Name, string(Terminated)).Set(metrics[Terminated])
		instancesGauge.WithLabelValues(scrapeTarget.Name, string(Failure)).Set(metrics[Failure])

		demand := scalingFunc(scrapeTarget, instances, metricValue)

		instanceDemandGauge.WithLabelValues(scrapeTarget.Name).Set(float64(demand.Result))
		metricValueGauge.WithLabelValues(scrapeTarget.Name).Set(metricValue)
		return demand
	}
}
