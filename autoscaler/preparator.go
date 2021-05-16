package autoscaler

import "github.com/cbrgm/cloudburst/cloudburst"

// preparator modifies instance states stored in the state to fulfill scaling needs calculated by a requestCalculator.
// Create preparator instances with newPreparator.
type preparator struct {
	state State
}

// newPreparator creates a new preparator. The provided State is used to access and modify instance states
// stored by a database provider implementation. The provided requestCalculator is used to calculate
func newPreparator(state State) *preparator {
	return &preparator{
		state: state,
	}
}

func (r *preparator) Prepare(result cloudburst.ScalingResult, scrapeTarget *cloudburst.ScrapeTarget) error {
	for _, value := range result.Result {
		err := r.prepareForProvider(value, scrapeTarget)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *preparator) prepareForProvider(value *cloudburst.ResultValue, scrapeTarget *cloudburst.ScrapeTarget) error {
	// if demand is in range of threshold, we don't request/suspend new instances

	sumTerminating := cloudburst.CountActiveInstances(cloudburst.GetInstancesByStatus(value.Instances, cloudburst.Running), false)
	sumProgressActive := cloudburst.CountActiveInstances(cloudburst.GetInstancesByStatus(value.Instances, cloudburst.Progress), true)
	sumRunning := cloudburst.CountActiveInstances(cloudburst.GetInstancesByStatus(value.Instances, cloudburst.Running), true)

	current := (sumRunning - sumTerminating) + sumProgressActive
	demand := value.InstanceDemand

	delta := demand - current

	//if r.threshold.inRange(delta) {
	//	delta = 0
	//}

	if delta == 0 {
		return r.thresholdEquals(scrapeTarget, value.Instances)
	}

	if delta > 0 {
		return r.thresholdAbove(scrapeTarget, value.Instances, value.Provider, delta)
	}

	if delta < 0 {
		return r.thresholdBelow(scrapeTarget, value.Instances, delta)
	}
	return nil
}

func (r *preparator) thresholdEquals(scrapeTarget *cloudburst.ScrapeTarget, instances []*cloudburst.Instance) error {
	pendingInstances := cloudburst.GetInstancesByStatus(instances, cloudburst.Pending)
	err := r.state.RemoveInstances(scrapeTarget.Name, pendingInstances)
	if err != nil {
		return err
	}
	return nil
}

func (r *preparator) thresholdAbove(scrapeTarget *cloudburst.ScrapeTarget, instances []*cloudburst.Instance, provider string, demand int) error {
	pendingInstances := cloudburst.GetInstancesByStatus(instances, cloudburst.Pending)

	var result = demand - len(pendingInstances)

	if result > 0 {
		// on a positive vars, create pending instances until we satisfy vars
		var toCreate []*cloudburst.Instance
		for i := 0; i < result; i++ {
			toCreate = append(toCreate, cloudburst.NewInstance(provider, scrapeTarget.InstanceSpec))
		}
		_, err := r.state.SaveInstances(scrapeTarget.Name, toCreate)
		if err != nil {
			return err
		}
	}
	if result < 0 {
		// on negative vars, convert vars to a positive int and delete the first n pending instances from the state.
		index := 0 - (result)
		toDelete := pendingInstances[:index]
		err := r.state.RemoveInstances(scrapeTarget.Name, toDelete)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *preparator) thresholdBelow(scrapeTarget *cloudburst.ScrapeTarget, instances []*cloudburst.Instance, demand int) error {

	// delete all pending instances first
	pendingInstances := cloudburst.GetInstancesByStatus(instances, cloudburst.Pending)
	err := r.state.RemoveInstances(scrapeTarget.Name, pendingInstances)
	if err != nil {
		return err
	}

	runningInstances := cloudburst.GetInstancesByStatus(instances, cloudburst.Running)
	activeInstances := cloudburst.GetActiveInstances(runningInstances, true)

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

func markToBeTerminated(instances []*cloudburst.Instance) []*cloudburst.Instance {
	var res = []*cloudburst.Instance{}
	for _, instance := range instances {
		instance.Active = false
		res = append(res, instance)
	}
	return res
}
