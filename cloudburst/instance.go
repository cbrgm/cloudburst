package cloudburst

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Status string

const (
	Unknown    Status = "unknown"
	Pending    Status = "pending"
	Running    Status = "running"
	Failure    Status = "failure"
	Progress   Status = "progress"
	Terminated Status = "terminated"
)

type (
	Instance struct {
		Name      string         `json:"name"`
		Endpoint  string         `json:"endpoint"`
		Active    bool           `json:"active"`
		Container ContainerSpec  `json:"container"`
		Status    InstanceStatus `json:"status"`
	}

	InstanceStatus struct {
		Agent   string    `json:"agent"`
		Status  Status    `json:"status"`
		Started time.Time `json:"started"`
	}
)

// TODO: add instance spec as attribute to Instances
func NewInstance(spec InstanceSpec) *Instance {
	return &Instance{
		Name:      newInstanceName(spec),
		Endpoint:  "",
		Active:    true,
		Container: spec.Container,
		Status: InstanceStatus{
			Agent:   "",
			Status:  Pending,
			Started: time.Now(),
		},
	}
}

func newInstanceName(spec InstanceSpec) string {
	return fmt.Sprintf("%s-%s-%s", spec.Container.Name, "instance", strconv.Itoa(rand.Intn(100000)))
}

func CountActiveInstances(instances []*Instance, active bool) int {
	var sum int
	for _, instance := range instances {
		if instance.Active == active {
			sum++
		}
	}
	return sum
}

func CountInstancesByStatus(instances []*Instance, status Status) int {
	var sum int
	for _, instance := range instances {
		if isMatchingStatus(instance, status) {
			sum++
		}
	}
	return sum
}

func GetActiveInstances(instances []*Instance, active bool) []*Instance {
	var res []*Instance
	for _, instance := range instances {
		if instance.Active == active {
			res = append(res, instance)
		}
	}
	return res
}

func GetInstancesByStatus(instances []*Instance, status Status) []*Instance {
	var res []*Instance
	for _, instance := range instances {
		if isMatchingStatus(instance, status) {
			res = append(res, instance)
		}
	}
	return res
}

func GetInstancesWithoutStatus(instances []*Instance, status Status) []*Instance {
	var res []*Instance
	for _, instance := range instances {
		if !isMatchingStatus(instance, status) {
			res = append(res, instance)
		}
	}
	return res
}

func isMatchingStatus(instance *Instance, status Status) bool {
	return instance.Status.Status == status
}
