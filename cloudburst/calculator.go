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
	sumEffective := math.Ceil((metricValue)*(sumInternal+sumExternal))

	demandTotal := sumEffective - sumInternal

	println("Values:")
	println(sumInternal)
	println(sumExternal)
	println(demandTotal)
	println("")

	res := ScalingResult{
		Result: []*ResultValue{},
	}

	for provider, weight := range scrapeTarget.ProviderSpec.Weights {
		val := calculateDemandForProvider(instances, demandTotal, provider, weight)
		res.Result = append(res.Result, &val)
	}

	handleDifference(res, demandTotal)

	for _, val := range res.Result {
		println(val.Provider)
		println(val.InstanceDemand)
		println(val.Weight)
		println("")
	}
	return res
}

func calculateDemandForProvider(instances []*Instance, demandTotal float64, provider string, weight float32) ResultValue {
	in := GetInstancesByProvider(instances, provider)
	demand := math.Round(demandTotal * float64(weight))
	return ResultValue{
		Provider:       provider,
		Weight:         weight,
		InstanceDemand: int(demand),
		Instances:      in,
	}
}

func handleDifference(res ScalingResult, demandTotal float64) ScalingResult {
	// sort provider weight ascending
	ByResultValues(func(p1, p2 *ResultValue) bool {
		return p1.Weight > p2.Weight
	}).Sort(res.Result)

	// calculate diff
	diff := demandTotal - float64(res.Sum())

	// handle diff
	if diff == 0 {
		return res
	}
	if diff > 0 {
		res.Result[0].InstanceDemand = res.Result[0].InstanceDemand + int(diff)
	}
	if diff < 0 {
		handleNegativeDifference(res, diff)
	}
	return res
}

func handleNegativeDifference(res ScalingResult, diff float64) ScalingResult {
	for _, in := range res.Result {
		if in.InstanceDemand+int(diff) <= 0 {
			diff = diff + float64(in.InstanceDemand)
			in.InstanceDemand = 0
			continue
		} else {
			in.InstanceDemand = in.InstanceDemand + int(diff)
			break
		}
	}
	return res
}
