package cloudburst

type (
	ScrapeTarget struct {
		Name         string       `json:"name"`
		Path         string       `json:"path"`
		Description  string       `json:"description"`
		Query        string       `json:"query"`
		InstanceSpec InstanceSpec `json:"instanceSpec"`
		StaticSpec   StaticSpec   `json:"static"`
	}

	InstanceSpec struct {
		Container ContainerSpec `json:"container"`
	}

	ContainerSpec struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	StaticSpec struct {
		Endpoints []string `json:"endpoints"`
	}
)
