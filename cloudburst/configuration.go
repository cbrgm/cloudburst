package cloudburst

type Configuration struct {
	PrometheusURL string `json:"prometheus_url"`
	ScrapeTargets []struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		Query        string `json:"query"`
		InstanceSpec struct {
			ContainerSpec struct {
				Name  string `json:"name"`
				Image string `json:"image"`
			} `json:"container"`
		} `json:"spec"`
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
			Query:       item.Query,
			InstanceSpec: InstanceSpec{
				Container: ContainerSpec{
					Name:  item.InstanceSpec.ContainerSpec.Name,
					Image: item.InstanceSpec.ContainerSpec.Image,
				},
			},
		})
	}

	return scrapeTargets, nil
}

func validateConfig(config Configuration) error {
	// TODO: add config validation checks
	return nil
}
