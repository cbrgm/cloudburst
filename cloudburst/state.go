package cloudburst

type StateProvider interface {
	ListScrapeTargets() ([]ScrapeTarget, error)
	GetInstance(name string) (Instance, error)
	GetInstancesForTarget(scrapeTarget string) ([]Instance, error)
	RemoveInstances(instances []Instance) error
	RemoveInstance(instance Instance) error
	SaveInstances(instances []Instance) ([]Instance, error)
	SaveInstance(instance Instance) (Instance, error)
}

type State struct {
	stateProvider StateProvider
}

func NewStateWithProvider(provider StateProvider) *State {
	return &State{
		stateProvider: provider,
	}
}

func (state *State) ListScrapeTargets() ([]ScrapeTarget, error) {
	return state.stateProvider.ListScrapeTargets()
}

func (state *State) GetInstance(name string) (Instance, error) {
	return state.stateProvider.GetInstance(name)
}

func (state *State) GetInstancesForTarget(scrapeTarget string) ([]Instance, error) {
	return state.stateProvider.GetInstancesForTarget(scrapeTarget)
}

func (state *State) RemoveInstances(instances []Instance) error {
	return state.stateProvider.RemoveInstances(instances)
}

func (state *State) RemoveInstance(instance Instance) error {
	return state.stateProvider.RemoveInstance(instance)
}

func (state *State) SaveInstances(scrapeTarget string, instances []Instance) ([]Instance, error) {
	return state.stateProvider.SaveInstances(instances)
}

func (state *State) SaveInstance(scrapeTarget string, instance Instance) (Instance, error) {
	return state.stateProvider.SaveInstance(instance)
}
