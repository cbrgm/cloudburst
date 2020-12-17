package cloudburst

import (
	"sync"
)

type EventType string

const (
	InstanceSaveEvent   EventType = "save"
	InstanceRemoveEvent EventType = "remove"
)

// InstanceEvent is an event to be pushed to the instance event bus.
type InstanceEvent struct {
	EventType    EventType `json:"type"`

	ScrapeTarget string    `json:"target"`
	Instance     *Instance `json:"data"`
}

// NewInstanceEvent creates a new event to be published via instance event bus.
func NewInstanceEvent(eventType EventType, instance *Instance) InstanceEvent {
	return InstanceEvent{
		EventType: eventType,
		Instance:  instance,
	}
}

// Events represents an event bus to broadcast instance updates
type Events struct {
	instanceSubscribers map[int64]chan InstanceEvent
	instanceLock        sync.RWMutex
}

// NewEvents creates a new event bus for instance events
func NewEvents() *Events {
	return &Events{
		instanceSubscribers: make(map[int64]chan InstanceEvent),
	}
}

// Subscription represents a subscription to the Events channel.
// It is used to identify a registered channel
type Subscription struct {
	id int64
}

func (e *Events) SubscribeToInstanceEvents(channel chan InstanceEvent) Subscription {
	e.instanceLock.Lock()
	defer e.instanceLock.Unlock()

	id := int64((len(e.instanceSubscribers)) + 1)
	e.instanceSubscribers[id] = channel
	return Subscription{id: id}
}

func (e *Events) UnsubscribeFromInstanceEvents(s Subscription) {
	e.instanceLock.Lock()
	defer e.instanceLock.Unlock()

	delete(e.instanceSubscribers, s.id)
}

func (e *Events) PublishInstanceEvent(event InstanceEvent) {
	e.instanceLock.RLock()
	defer e.instanceLock.RUnlock()

	for _, subscriber := range e.instanceSubscribers {
		subscriber <- event
	}
}
