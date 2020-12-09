package state

import "github.com/cbrgm/cloudburst/cloudburst"

type StateProvider interface {
	ListScrapeTargets() ([]cloudburst.ScrapeTarget, error)
	GetInstances(scrapeTarget string) []cloudburst.Instance
	GetInstance(scrapeTarget string) []cloudburst.Instance
	UpdateInstances(scrapeTarget string, instances []cloudburst.Instance) ([]cloudburst.Instance, error)
	RemoveInstances(scrapeTarget string, instances []cloudburst.Instance) error
	RemoveInstance(scrapeTarget string, instance cloudburst.Instance) error
	CreateInstances(scrapeTarget string, instances []cloudburst.Instance) ([]cloudburst.Instance, error)
	CreateInstance(scrapeTarget string, instance cloudburst.Instance) (cloudburst.Instance, error)
}

type State struct {
	StateProvider StateProvider
}

func NewStateWithProvider(provider StateProvider) *State {
	return &State{
		StateProvider: provider,
	}
}

func (state *State) ListScrapeTargets() ([]cloudburst.ScrapeTarget, error) {
	return state.StateProvider.ListScrapeTargets()
}

func (state *State) UpdateInstances(scrapeTarget string, instances []cloudburst.Instance) ([]cloudburst.Instance, error) {
	return state.StateProvider.UpdateInstances(scrapeTarget, instances)
}

func (state *State) RemoveInstances(scrapeTarget string, instances []cloudburst.Instance) error {
	return state.StateProvider.RemoveInstances(scrapeTarget, instances)
}

func (state *State) RemoveInstance(scrapeTarget string, instance cloudburst.Instance) error {
	return state.StateProvider.RemoveInstance(scrapeTarget, instance)
}

func (state *State) CreateInstances(scrapeTarget string, instances []cloudburst.Instance) ([]cloudburst.Instance, error) {
	return state.StateProvider.CreateInstances(scrapeTarget, instances)
}

func (state *State) CreateInstance(scrapeTarget string, instance cloudburst.Instance) (cloudburst.Instance, error) {
	return state.StateProvider.CreateInstance(scrapeTarget, instance)
}
