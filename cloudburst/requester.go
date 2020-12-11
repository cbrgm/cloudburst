package cloudburst

const (
	threshold = 0
)

// requester modifies instance states stored in the state to fulfill scaling needs calculated by a requestCalculator.
// Create requester instances with newRequester.
type requester struct {
	state State
}

// newRequester creates a new requester. The provided State is used to access and modify instance states
// stored by a database provider implementation. The provided requestCalculator is used to calculate
func newRequester(state State) *requester {
	return &requester{
		state: state,
	}
}

func (r *requester) ProcessDemand(demand instanceDemand, instances []Instance, spec InstanceSpec) error {
	result := demand.Result
	if result == threshold {
		return r.thresholdEquals(instances)
	}

	if result > threshold {
		return r.thresholdAbove(result, instances, spec)
	}

	if result < threshold {
		return r.thresholdBelow(result, instances, spec)
	}
	return nil
}

func (r *requester) thresholdEquals(instances []Instance) error {
	pendingInstances := GetInstancesByStatus(instances, Pending)
	err := r.state.RemoveInstances(pendingInstances)
	if err != nil {
		return err
	}

	return nil
}

func (r *requester) thresholdAbove(demand int, instances []Instance, spec InstanceSpec) error {
	// TODO: implement
	return nil
}

func (r *requester) thresholdBelow(demand int, instances []Instance, spec InstanceSpec) error {
	// TODO: implement
	return nil
}
