package cloudburst

import (
	"github.com/prometheus/client_golang/prometheus"
	"math"
)

// ScalingFunc calculates the demand for Instance objects.
//
// instances slice is a list of all instances  to be filtered for calculation.
// metricValue value is the result of a metric query.
// CalculateDemand returns ScalingResult.
type ScalingFunc func(scrapeTarget *ScrapeTarget, instances []*Instance, metricValue float64) ScalingResult

// ScalingResult represents the calculated demand for instances
type ScalingResult struct {
	Result []ResultValue // the rounded demand for new instances
}

type ResultValue struct {
	Provider       string
	InstanceDemand int
	Instances      []*Instance
}

// NewDefaultScalingFunc returns the default scaling func
func NewDefaultScalingFunc() ScalingFunc {
	return func(scrapeTarget *ScrapeTarget, instances []*Instance, metricValue float64) ScalingResult {

		sumInternal := float64(len(scrapeTarget.StaticSpec.Endpoints))
		sumExternal := float64(CountInstancesByStatus(instances, Running))
		sumEffective := math.Round((metricValue-1)*(sumInternal+sumExternal) + 0.5)

		res := ScalingResult{
			Result: []ResultValue{},
		}

		for provider, weight := range scrapeTarget.ProviderSpec.Weights {
			val := calculateDemandForProvider(instances, sumEffective, provider, weight)
			res.Result = append(res.Result, val)
		}
		return res
	}
}

func calculateDemandForProvider(instances []*Instance, sumEffective float64, provider string, weight float32) ResultValue {
	in := GetInstancesByProvider(instances, provider)
	sumTerminating := float64(CountActiveInstances(in, false))
	sumProgressActive := float64(CountActiveInstances(GetInstancesByStatus(in, Progress), true))

	demand := ((math.Round(sumEffective * float64(weight))) + sumTerminating) - sumProgressActive
	return ResultValue{
		Provider:       provider,
		InstanceDemand: int(demand),
		Instances:      in,
	}
}

func InstrumentedScalingFunc(r *prometheus.Registry, scalingFunc ScalingFunc) ScalingFunc {

	instancesGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cloudburst_api_instances_total",
		Help: "instances total",
	}, []string{"target", "provider", "status"})

	instanceDemandGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cloudburst_api_demand",
		Help: "calculated instance demand for the current scaling iteration",
	}, []string{"target", "provider"})

	metricValueGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cloudburst_api_query",
		Help: "metric value received from prometheus for the current scaling iteration",
	}, []string{"target"})

	r.MustRegister(instancesGauge)
	r.MustRegister(instanceDemandGauge)
	r.MustRegister(metricValueGauge)

	return func(scrapeTarget *ScrapeTarget, instances []*Instance, metricValue float64) ScalingResult {
		metricValueGauge.WithLabelValues(scrapeTarget.Name).Set(metricValue)
		for provider, _ := range scrapeTarget.ProviderSpec.Weights {
			in := GetInstancesByProvider(instances, provider)
			metrics := make(map[Status]float64)
			for _, instance := range in {
				if val, ok := metrics[instance.Status.Status]; ok {
					metrics[instance.Status.Status] = val + 1
				} else {
					metrics[instance.Status.Status] = 1
				}
			}
			instancesGauge.WithLabelValues(scrapeTarget.Name, provider, string(Pending)).Set(metrics[Pending])
			instancesGauge.WithLabelValues(scrapeTarget.Name, provider, string(Progress)).Set(metrics[Progress])
			instancesGauge.WithLabelValues(scrapeTarget.Name, provider, string(Running)).Set(metrics[Running])
			instancesGauge.WithLabelValues(scrapeTarget.Name, provider, string(Terminated)).Set(metrics[Terminated])
			instancesGauge.WithLabelValues(scrapeTarget.Name, provider, string(Failure)).Set(metrics[Failure])
		}

		result := scalingFunc(scrapeTarget, instances, metricValue)

		for _, result := range result.Result {
			instanceDemandGauge.WithLabelValues(scrapeTarget.Name, result.Provider).Set(float64(result.InstanceDemand))
		}
		return result
	}
}
