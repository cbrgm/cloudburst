package cloudburst

type State interface {
	ListScrapeTargets() ([]*ScrapeTarget, error)
	GetInstance(scrapeTarget string, name string) (*Instance, error)
	GetInstances(scrapeTarget string) ([]*Instance, error)
	RemoveInstances(scrapeTarget string, instances []*Instance) error
	RemoveInstance(scrapeTarget string, instance *Instance) error
	SaveInstances(scrapeTarget string, instances []*Instance) ([]*Instance, error)
	SaveInstance(scrapeTarget string, instance *Instance) (*Instance, error)
}