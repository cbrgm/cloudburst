package cloudburst

type Configuration struct {
	PrometheusURL string `json:"prometheus_url"`
	ScrapeTargets []struct {
		Name         string `json:"name"`
		Path         string `json:"path"`
		Description  string `json:"description"`
		Query        string `json:"query"`
		InstanceSpec struct {
			ContainerSpec struct {
				Name  string `json:"name"`
				Image string `json:"image"`
			} `json:"container"`
		} `json:"spec"`
		StaticSpec struct {
			Endpoints []string `json:"endpoints"`
		} `json:"static"`
	} `json:"targets"`
}

func ParseConfiguration(config Configuration) ([]*ScrapeTarget, error) {
	var scrapeTargets []*ScrapeTarget

	err := validateConfig(config)
	if err != nil {
		return scrapeTargets, err
	}

	for _, item := range config.ScrapeTargets {
		scrapeTargets = append(scrapeTargets, &ScrapeTarget{
			Name:        item.Name,
			Description: item.Description,
			Path:        item.Path,
			Query:       item.Query,
			InstanceSpec: InstanceSpec{
				Container: ContainerSpec{
					Name:  item.InstanceSpec.ContainerSpec.Name,
					Image: item.InstanceSpec.ContainerSpec.Image,
				},
			},
			StaticSpec: StaticSpec{
				Endpoints: item.StaticSpec.Endpoints,
			},
		})
	}

	return scrapeTargets, nil
}

func validateConfig(config Configuration) error {
	// TODO: add config validation checks
	return nil
}
