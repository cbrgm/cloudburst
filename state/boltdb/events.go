package boltdb

import "github.com/cbrgm/cloudburst/cloudburst"

type Events struct {
	*BoltDB
	events *cloudburst.Events
}

func NewEvents(db *BoltDB, events *cloudburst.Events) *Events {
	return &Events{
		BoltDB: db,
		events: events,
	}
}

func (e *Events) ListScrapeTargets() ([]*cloudburst.ScrapeTarget, error) {
	return e.BoltDB.ListScrapeTargets()
}

func (e *Events) GetInstance(scrapeTarget string, name string) (*cloudburst.Instance, error) {
	return e.BoltDB.GetInstance(scrapeTarget, name)
}

func (e *Events) GetInstances(scrapeTarget string) ([]*cloudburst.Instance, error) {
	return e.BoltDB.GetInstances(scrapeTarget)
}

func (e *Events) RemoveInstances(scrapeTarget string, instances []*cloudburst.Instance) error {
	if len(instances) == 0 {
		return nil
	}

	err := e.BoltDB.RemoveInstances(scrapeTarget, instances)
	if err != nil {
		return err
	}
	e.events.PublishInstanceEvent(cloudburst.InstanceEvent{
		EventType:    cloudburst.InstanceRemoveEvent,
		ScrapeTarget: scrapeTarget,
		Instances:    instances,
	})

	return err
}

func (e *Events) RemoveInstance(scrapeTarget string, instance *cloudburst.Instance) error {
	if instance == nil {
		return nil
	}

	err := e.BoltDB.RemoveInstance(scrapeTarget, instance)
	if err != nil {
		return err
	}

	e.events.PublishInstanceEvent(cloudburst.InstanceEvent{
		EventType:    cloudburst.InstanceRemoveEvent,
		ScrapeTarget: scrapeTarget,
		Instances:    []*cloudburst.Instance{instance},
	})
	return err
}

func (e *Events) SaveInstances(scrapeTarget string, instances []*cloudburst.Instance) ([]*cloudburst.Instance, error) {
	if len(instances) == 0 {
		return instances, nil
	}
	saved, err := e.BoltDB.SaveInstances(scrapeTarget, instances)
	if err != nil {
		return saved, err
	}

	e.events.PublishInstanceEvent(cloudburst.InstanceEvent{
		EventType:    cloudburst.InstanceSaveEvent,
		ScrapeTarget: scrapeTarget,
		Instances:    saved,
	})
	return saved, nil
}

func (e *Events) SaveInstance(scrapeTarget string, instance *cloudburst.Instance) (*cloudburst.Instance, error) {
	if instance == nil {
		return nil, nil
	}
	saved, err := e.BoltDB.SaveInstance(scrapeTarget, instance)
	if err != nil {
		return saved, err
	}

	e.events.PublishInstanceEvent(cloudburst.InstanceEvent{
		EventType:    cloudburst.InstanceSaveEvent,
		ScrapeTarget: scrapeTarget,
		Instances:    []*cloudburst.Instance{saved},
	})
	return saved, err
}
