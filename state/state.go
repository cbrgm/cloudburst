package state

import "github.com/cbrgm/cloudburst/cloudburst"

type StateProvider interface {
	ListScrapeTargets() ([]cloudburst.ScrapeTarget, error)
	GetInstance(name string) (cloudburst.Instance, error)
	GetInstancesForTarget(scrapeTarget string) ([]cloudburst.Instance, error)
	RemoveInstances(instances []cloudburst.Instance) error
	RemoveInstance(instance cloudburst.Instance) error
	SaveInstances(instances []cloudburst.Instance) ([]cloudburst.Instance, error)
	SaveInstance(instance cloudburst.Instance) (cloudburst.Instance, error)
}

type State struct {
	stateProvider StateProvider
}

func NewStateWithProvider(provider StateProvider) *State {
	return &State{
		stateProvider: provider,
	}
}

func (state *State) ListScrapeTargets() ([]cloudburst.ScrapeTarget, error) {
	return state.stateProvider.ListScrapeTargets()
}

func (state *State) GetInstance(name string) (cloudburst.Instance, error) {
	return state.stateProvider.GetInstance(name)
}

func (state *State) GetInstancesForTarget(scrapeTarget string) ([]cloudburst.Instance, error) {
	return state.stateProvider.GetInstancesForTarget(scrapeTarget)
}

func (state *State) RemoveInstances(instances []cloudburst.Instance) error {
	return state.stateProvider.RemoveInstances(instances)
}

func (state *State) RemoveInstance(instance cloudburst.Instance) error {
	return state.stateProvider.RemoveInstance(instance)
}

func (state *State) SaveInstances(scrapeTarget string, instances []cloudburst.Instance) ([]cloudburst.Instance, error) {
	return state.stateProvider.SaveInstances(instances)
}

func (state *State) SaveInstance(scrapeTarget string, instance cloudburst.Instance) (cloudburst.Instance, error) {
	return state.stateProvider.SaveInstance(instance)
}
