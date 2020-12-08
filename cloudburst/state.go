package cloudburst

type StateProvider interface {
	ListScrapeTargets() ([]ScrapeTarget, error)
	UpdateInstances(scrapeTarget string, instances []Instance) ([]Instance, error)
	RemoveInstances(scrapeTarget string, instances []Instance) error
	RemoveInstance(scrapeTarget string, instance Instance) error
	CreateInstances(scrapeTarget string, instances []Instance) ([]Instance, error)
	CreateInstance(scrapeTarget string, instance Instance) (Instance, error)
}

type State struct {
	StateProvider StateProvider
}

func (state *State) ListScrapeTargets() ([]ScrapeTarget, error) {
	return state.StateProvider.ListScrapeTargets()
}

func (state *State) UpdateInstances(scrapeTarget string, instances []Instance) ([]Instance, error) {
	return state.StateProvider.UpdateInstances(scrapeTarget, instances)
}

func (state *State) RemoveInstances(scrapeTarget string, instances []Instance) error {
	return state.StateProvider.RemoveInstances(scrapeTarget, instances)
}

func (state *State) RemoveInstance(scrapeTarget string, instance Instance) error {
	return state.StateProvider.RemoveInstance(scrapeTarget, instance)
}

func (state *State) CreateInstances(scrapeTarget string, instances []Instance) ([]Instance, error) {
	return state.StateProvider.CreateInstances(scrapeTarget, instances)
}

func (state *State) CreateInstance(scrapeTarget string, instance Instance) (Instance, error) {
	return state.StateProvider.CreateInstance(scrapeTarget, instance)
}

func NewVolatileState(initialState []ScrapeTarget) (*State, error) {
	provider := newVolatileState(initialState)

	return &State{
		StateProvider: provider,
	}, nil
}

type volatileState struct {
	State map[string]ScrapeTarget
}

func newVolatileState(scrapeTargets []ScrapeTarget) *volatileState {
	state := map[string]ScrapeTarget{}
	for _, target := range scrapeTargets {
		state[target.Name] = target
	}
	return &volatileState{State: state}
}

func (v *volatileState) ListScrapeTargets() ([]ScrapeTarget, error) {
	var res []ScrapeTarget
	for _, v := range v.State {
		res = append(res, v)
	}
	return res, nil
}
func (v *volatileState) UpdateInstances(scrapeTarget string, instances []Instance) ([]Instance, error) {
	return nil, nil
}

func (v *volatileState) RemoveInstances(scrapeTarget string, instances []Instance) error {
	return nil
}
func (v *volatileState) RemoveInstance(scrapeTarget string, instance Instance) error {
	target := v.State[scrapeTarget]
	for index, item := range target.Instances {
		if item.Name == instance.Name {
			target.Instances = append(target.Instances[:index],target.Instances[index+1:]...)
			break
		}
	}
	v.State[scrapeTarget] = target
	return nil
}
func (v *volatileState) CreateInstances(scrapeTarget string, instances []Instance) ([]Instance, error) {
	for _, instance := range instances {
		_, err := v.CreateInstance(scrapeTarget, instance)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
func (v *volatileState) CreateInstance(scrapeTarget string, instance Instance) (Instance, error) {
	target := v.State[scrapeTarget]

	target.Instances = append(target.Instances, instance)
	v.State[scrapeTarget] = target
	return instance, nil
}
