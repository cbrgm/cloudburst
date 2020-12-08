package cloudburst

import "time"

type (
	Instance struct {
		Name     string         `json:"name"`
		Endpoint string         `json:"endpoint"`
		Active   bool           `json:"active"`
		Status   InstanceStatus `json:"status"`
	}

	InstanceStatus struct {
		Agent   string    `json:"agent"`
		Status  string    `json:"status"`
		Started time.Time `json:"started"`
	}
)
