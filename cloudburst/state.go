package cloudburst

type StateProvider interface {
	ListScrapeTargets() ([]ScrapeTarget, error)
	UpdateInstances(targetName string, instances []Instance) ([]Instance, error)
}

type State struct {
	StateProvider StateProvider
}

func (db *State) ListScrapeTargets() ([]ScrapeTarget, error) {
	return db.StateProvider.ListScrapeTargets()
}

func (db *State) UpdateInstances(targetName string, instances []Instance) ([]Instance, error) {
	return db.StateProvider.UpdateInstances(targetName, instances)
}

func NewVolatileState(initialState []ScrapeTarget) (*State, error) {
	provider := newVolatileState(initialState)

	return &State{
		StateProvider: provider,
	}, nil
}

type volatileState struct {
	State []ScrapeTarget
}

func newVolatileState(scrapeTargets []ScrapeTarget) *volatileState {
	return &volatileState{State: scrapeTargets}
}

func (v *volatileState) ListScrapeTargets() ([]ScrapeTarget, error) {
	return v.State, nil
}

func (v *volatileState) UpdateInstances(targetName string, instances []Instance) ([]Instance, error) {
	return nil, nil
}
