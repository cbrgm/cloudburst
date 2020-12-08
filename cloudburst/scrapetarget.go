package cloudburst

type (
	ScrapeTarget struct {
		Name         string       `json:"name"`
		Description  string       `json:"description"`
		Query        string       `json:"query"`
		Value        float64      `json:"value"`
		InstanceSpec InstanceSpec `json:"instanceSpec"`
		Instances    []Instance   `json:"instances"`
	}

	InstanceSpec struct {
		Container ContainerSpec `json:"container"`
	}

	ContainerSpec struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}
)
