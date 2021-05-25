package main

import (
	"context"
	"github.com/cbrgm/cloudburst/cloudburst"
)

type InstanceLifecycleHandler interface {
	// CreateInstance takes an instance spec and deploys it within the cloud provider.
	CreateInstance(ctx context.Context, inst *cloudburst.Instance) error
	// UpdateInstance takes an instance and updates it within the cloud provider.
	UpdateInstance(ctx context.Context, inst *cloudburst.Instance) error
	// DeleteInstance takes an instance and deletes it from the cloud provider.
	DeleteInstance(ctx context.Context, inst *cloudburst.Instance) error
	// GetInstance retrieves an instance by name from the provider (can be cached).
	// The instance returned is expected to be immutable, and may be accessed
	// concurrently outside of the calling goroutine.
	GetInstance(ctx context.Context, namespace, name string) (*cloudburst.Instance, error)
	// GetInstanceStatus retrieves the status of an instance by name from the provider.
	// The InstanceStatus returned is expected to be immutable, and may be accessed
	// concurrently outside of the calling goroutine.
	GetInstanceStatus(ctx context.Context, namespace, name string) (*cloudburst.InstanceStatus, error)
	// GetInstances retrieves a list of all instances running on the cloud provider (can be cached).
	// The Instances returned are expected to be immutable, and may be accessed
	// concurrently outside of the calling goroutine.
	GetInstances(context.Context) ([]*cloudburst.Instance, error)
}