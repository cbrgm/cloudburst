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
		Active   bool           `json:"active"`
		Status   InstanceStatus `json:"status"`
	}

	InstanceStatus struct {
		Agent   string    `json:"agent"`
		Status  Status    `json:"status"`
		Started time.Time `json:"started"`
	}
)
