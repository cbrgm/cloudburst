package cloudburst

type AutoScaler struct {
	state     State
	requester *requester
}

func NewAutoScaler(state State) *AutoScaler {
	return &AutoScaler{
		state:      state,
		requester:  newRequester(state),
	}
}

func (scaling *AutoScaler) Scale(scrapeTarget ScrapeTarget, queryResult float64) error {
	instances, err := scaling.state.GetInstancesForTarget(scrapeTarget.Name)
	if err != nil {
		return err
	}

	spec := scrapeTarget.InstanceSpec
	demand := scaling.calculateDemand(instances, queryResult)

	err = scaling.processDemand(demand, instances, spec)
	if err != nil {
		return err
	}

	return nil
}

type instanceDemand struct {
	Result int
}

// CalculateDemand calculates the demand for Instance objects. The provided instances slice is a list of all instances
// to be filtered for calculation. The provided queryResult value is the result of a metric query.
// CalculateDemand returns instanceDemand.
func (scaling *AutoScaler) calculateDemand(instances []Instance, queryResult float64) instanceDemand {
	sumTerminating := CountTerminatingInstances(instances)
	sumProgress := CountInstancesByStatus(instances, Progress)

	var demand = (int(queryResult) + sumTerminating) - sumProgress

	return instanceDemand{
		Result: demand,
	}
}

func (scaling *AutoScaler) processDemand(demand instanceDemand, instances []Instance, spec InstanceSpec) error {
	return scaling.requester.ProcessDemand(demand, instances, spec)
}
