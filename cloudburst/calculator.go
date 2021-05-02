package cloudburst

import "math"

type ReactiveThreshold struct {
	threshold Threshold
}

func NewReactiveThreshold() *ReactiveThreshold {
	return &ReactiveThreshold{
		threshold: newThreshold(0, 0),
	}
}

func NewReactiveThresholdWithThreshold(threshold Threshold) *ReactiveThreshold {
	return &ReactiveThreshold{
		threshold: threshold,
	}
}

func (c *ReactiveThreshold) Calculate(scrapeTarget *ScrapeTarget, instances []*Instance, metricValue float64) ScalingResult {
	sumInternal := float64(len(scrapeTarget.StaticSpec.Endpoints))
	sumExternal := float64(CountInstancesByStatus(instances, Running))
	sumEffective := math.Round((metricValue-1)*(sumInternal+sumExternal) + 0.5)

	res := ScalingResult{
		Result: []*ResultValue{},
	}

	for provider, weight := range scrapeTarget.ProviderSpec.Weights {
		val := calculateDemandForProvider(instances, sumEffective, provider, weight)
		res.Result = append(res.Result, &val)
	}

	ByResultValues(func(p1, p2 *ResultValue) bool {
		return p1.Weight > p2.Weight
	}).Sort(res.Result)

	demand := sumEffective - float64(res.Sum())
	if demand == 0 {
		return res
	}
	if demand > 0 {
		res.Result[0].InstanceDemand = res.Result[0].InstanceDemand + int(demand)
	}
	if demand < 0 {
		for _, in := range res.Result {
			if in.InstanceDemand-int(demand) <= 0 {
				continue
			} else {
				in.InstanceDemand = in.InstanceDemand + int(demand)
				break
			}
		}
	}

	// add threshold
	return res
}

func calculateDemandForProvider(instances []*Instance, sumEffective float64, provider string, weight float32) ResultValue {
	in := GetInstancesByProvider(instances, provider)
	sumTerminating := float64(CountActiveInstances(in, false))
	sumProgressActive := float64(CountActiveInstances(GetInstancesByStatus(in, Progress), true))

	demand := ((math.Round(sumEffective * float64(weight))) + sumTerminating) - sumProgressActive
	return ResultValue{
		Provider:       provider,
		Weight:         weight,
		InstanceDemand: int(demand),
		Instances:      in,
	}
}
