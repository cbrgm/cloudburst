package cloudburst

type StateProvider interface {
	ListScrapeTargets() ([]ScrapeTarget, error)
	GetInstance(scrapeTarget string, name string) (Instance, error)
	GetInstances(scrapeTarget string) ([]Instance, error)
	RemoveInstances(scrapeTarget string, instances []Instance) error
	RemoveInstance(scrapeTarget string, instance Instance) error
	SaveInstances(scrapeTarget string, instances []Instance) ([]Instance, error)
	SaveInstance(scrapeTarget string, instance Instance) (Instance, error)
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

func (state *State) GetInstance(scrapeTarget string, name string) (Instance, error) {
	return state.stateProvider.GetInstance(scrapeTarget, name)
}

func (state *State) GetInstances(scrapeTarget string) ([]Instance, error) {
	return state.stateProvider.GetInstances(scrapeTarget)
}

func (state *State) RemoveInstances(scrapeTarget string, instances []Instance) error {
	return state.stateProvider.RemoveInstances(scrapeTarget, instances)
}

func (state *State) RemoveInstance(scrapeTarget string, instance Instance) error {
	return state.stateProvider.RemoveInstance(scrapeTarget, instance)
}

func (state *State) SaveInstances(scrapeTarget string, instances []Instance) ([]Instance, error) {
	return state.stateProvider.SaveInstances(scrapeTarget, instances)
}

func (state *State) SaveInstance(scrapeTarget string, instance Instance) (Instance, error) {
	return state.stateProvider.SaveInstance(scrapeTarget, instance)
}
