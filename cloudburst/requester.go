package cloudburst

type Threshold struct {
	Upper int
	Lower int
}

func (t *Threshold) inRange(i int) bool {
	return i >= t.Lower && i <= t.Upper
}

func newThreshold(upper, lower int) Threshold {
	return Threshold{
		Upper: upper,
		Lower: lower,
	}
}

// requester modifies instance states stored in the state to fulfill scaling needs calculated by a requestCalculator.
// Create requester instances with newRequester.
type requester struct {
	threshold Threshold
	state     State
}

// newRequester creates a new requester. The provided State is used to access and modify instance states
// stored by a database provider implementation. The provided requestCalculator is used to calculate
func newRequester(state State, threshold Threshold) *requester {
	return &requester{
		threshold: threshold,
		state:     state,
	}
}

func (r *requester) ProcessDemand(demand instanceDemand, instances []*Instance, scrapeTarget *ScrapeTarget) error {
	result := demand.Result

	// if demand is in range of threshold, we don't request/supend new instances
	if r.threshold.inRange(result) {
		result = 0
	}

	if result == 0 {
		return r.thresholdEquals(scrapeTarget, instances)
	}

	if result > 0 {
		return r.thresholdAbove(scrapeTarget, instances, result)
	}

	if result < 0 {
		return r.thresholdBelow(scrapeTarget, instances, result)
	}
	return nil
}

func (r *requester) thresholdEquals(scrapeTarget *ScrapeTarget, instances []*Instance) error {
	pendingInstances := GetInstancesByStatus(instances, Pending)
	err := r.state.RemoveInstances(scrapeTarget.Name, pendingInstances)
	if err != nil {
		return err
	}
	return nil
}

func (r *requester) thresholdAbove(scrapeTarget *ScrapeTarget, instances []*Instance, demand int) error {
	pendingInstances := GetInstancesByStatus(instances, Pending)

	var result = demand - len(pendingInstances)

	if result > 0 {
		// on a positive result, create pending instances until we satisfy result
		var toCreate []*Instance
		for i := 0; i < result; i++ {
			toCreate = append(toCreate, NewInstance(scrapeTarget.InstanceSpec))
		}
		_, err := r.state.SaveInstances(scrapeTarget.Name, toCreate)
		if err != nil {
			return err
		}
	}
	if result < 0 {
		// on negative result, convert result to a positive int and delete the first n pending instances from the state.
		index := 0 - (result)
		toDelete := pendingInstances[:index]
		err := r.state.RemoveInstances(scrapeTarget.Name, toDelete)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *requester) thresholdBelow(scrapeTarget *ScrapeTarget, instances []*Instance, demand int) error {

	// delete all pending instances first
	pendingInstances := GetInstancesByStatus(instances, Pending)
	err := r.state.RemoveInstances(scrapeTarget.Name, pendingInstances)
	if err != nil {
		return err
	}

	runningInstances := GetInstancesByStatus(instances, Running)
	activeInstances := GetActiveInstances(runningInstances, true)

	// we dont want a negative number here
	numToBeTerminated := 0 - (demand)

	if numToBeTerminated >= len(activeInstances) {
		// set all instances to active == false
		res := markToBeTerminated(activeInstances)
		_, err := r.state.SaveInstances(scrapeTarget.Name, res)
		if err != nil {
			return err
		}
	} else {
		// terminate instances according to numToBeTerminated
		toBeTerminated := activeInstances[:numToBeTerminated]
		res := markToBeTerminated(toBeTerminated)
		_, err := r.state.SaveInstances(scrapeTarget.Name, res)
		if err != nil {
			return err
		}
	}
	return nil
}

func markToBeTerminated(instances []*Instance) []*Instance {
	var res = []*Instance{}
	for _, instance := range instances {
		instance.Active = false
		res = append(res, instance)
	}
	return res
}
