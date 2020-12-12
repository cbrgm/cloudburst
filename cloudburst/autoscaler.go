package cloudburst

type AutoScaler struct {
	state     *State
	requester *requester
}

func NewAutoScaler(state *State) *AutoScaler {
	return &AutoScaler{
		state:     state,
		requester: newRequester(state),
	}
}

func (s *AutoScaler) Scale(scrapeTarget ScrapeTarget, queryResult float64) error {
	instances, err := s.state.GetInstances(scrapeTarget.Name)
	if err != nil {
		return err
	}

	demand := s.calculateDemand(instances, queryResult)
	err = s.processDemand(demand, instances, scrapeTarget)
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
func (s *AutoScaler) calculateDemand(instances []Instance, queryResult float64) instanceDemand {
	sumTerminating := CountInstancesByActiveStatus(instances, false)
	sumProgress := CountInstancesByStatus(instances, Progress)

	var demand = (int(queryResult) + sumTerminating) - sumProgress

	return instanceDemand{
		Result: demand,
	}
}

func (s *AutoScaler) processDemand(demand instanceDemand, instances []Instance, scrapeTarget ScrapeTarget) error {
	return s.requester.ProcessDemand(demand, instances, scrapeTarget)
}
