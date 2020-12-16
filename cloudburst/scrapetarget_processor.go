package cloudburst

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"time"
)

type ScrapeTargetProcessor struct {
	state      *State
	autoscaler *AutoScaler
}

func NewScrapeTargetProcessor(state *State) *ScrapeTargetProcessor {
	return &ScrapeTargetProcessor{
		state:      state,
		autoscaler: NewAutoScaler(state),
	}
}

func (sp *ScrapeTargetProcessor) ProcessScrapeTargets(prometheusURL string) error {
	client, err := api.NewClient(api.Config{Address: prometheusURL})
	if err != nil {
		return fmt.Errorf("failed to create Prometheus client: %w", err)
	}
	promAPI := prometheusv1.NewAPI(client)

	scrapeTargets, err := sp.state.ListScrapeTargets()
	if err != nil {
		return err
	}

	for _, target := range scrapeTargets {
		err = sp.processScrapeTarget(promAPI, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sp *ScrapeTargetProcessor) processScrapeTarget(promAPI prometheusv1.API, scrapeTarget *ScrapeTarget) error {
	value, _, err := promAPI.Query(context.TODO(), scrapeTarget.Query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to run processScrapeTargets: %w", err)
	}
	vec := value.(model.Vector)
	var queryResult float64

	for _, v := range vec {
		if v.Value.String() == "NaN" {
			queryResult = 0
		} else {
			queryResult = float64(v.Value)
		}
	}

	err = sp.autoscaler.Scale(scrapeTarget, queryResult)
	if err != nil {
		return err
	}

	return nil
}
