package cloudburst

type Configuration struct {
	PrometheusURL string `json:"prometheus_url"`
	ScrapeTargets []struct {
		Name         string `json:"name"`
		Path         string `json:"path"`
		Description  string `json:"description"`
		Query        string `json:"query"`
		ProviderSpec struct {
			Weights map[string]int `json:"weights"`
		} `json:"provider"`
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
			Name:         item.Name,
			Description:  item.Description,
			Path:         item.Path,
			Query:        item.Query,
			ProviderSpec: ProviderSpec{
				Weights: calculateWeights(item.ProviderSpec.Weights),
			},
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

func calculateWeights(weights map[string]int) map[string]float32 {
	var res = make(map[string]float32)
	var sumWeights int
	for _, v := range weights {
		sumWeights += v
	}

	for k, v := range weights {
		res[k] = float32(v) / float32(sumWeights)
	}
	return res
}

func validateConfig(config Configuration) error {
	return nil
}

/*
func validateScrapeTarget(target) error {
	validateProviderSpec()
	return nil
}

func validateProviderSpec(config Configuration) error {
	hasIllegalProviderWeights(config.)
}

func hasIllegalProviderWeights(weights map[string]int) error {
	var sumWeights int
	for _,v := range weights {
		if v < 0 {
			return errors.New("negative values not allowed")
		}
		sumWeights += v
	}
	if sumWeights <= 0 {
		return errors.New("provider weight sum must be greater then 0")
	}
}
*/
