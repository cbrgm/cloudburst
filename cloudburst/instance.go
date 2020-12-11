package cloudburst

import "time"

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
		Name     string         `json:"name"`
		Endpoint string         `json:"endpoint"`
		Target   string         `json:"target"`
		Active   bool           `json:"active"`
		Status   InstanceStatus `json:"status"`
	}

	InstanceStatus struct {
		Agent   string    `json:"agent"`
		Status  Status    `json:"status"`
		Started time.Time `json:"started"`
	}
)

// TODO: add instance spec as attribute to Instances
func NewInstance(spec InstanceSpec) Instance {
	return Instance{
		Name:     "foo",
		Endpoint: "",
		Target:   "",
		Active:   false,
		Status: InstanceStatus{
			Agent:   "",
			Status:  Pending,
			Started: time.Now(),
		},
	}
}

func CountTerminatingInstances(instances []Instance) int {
	var sum int
	for _, instance := range instances {
		if instance.Active == false {
			sum++
		}
	}
	return sum
}

func CountInstancesByStatus(instances []Instance, status Status) int {
	var sum int
	for _, instance := range instances {
		if isMatchingStatus(instance, status) {
			sum++
		}
	}
	return sum
}



func GetInstancesByStatus(instances []Instance, status Status) []Instance {
	var res []Instance
	for _, instance := range instances {
		if isMatchingStatus(instance, status) {
			res = append(res, instance)
		}
	}
	return res
}

func isMatchingStatus(instance Instance, status Status) bool {
	return instance.Status.Status == status
}
