package main

import (
	"context"
	"fmt"
	"github.com/cbrgm/cloudburst/cloudburst"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"time"
)

func processScrapeTargets(promAPI prometheusv1.API, state *cloudburst.State) error {
	scrapeTargets, err := state.ListScrapeTargets()
	if err != nil {
		return err
	}

	for _, target := range scrapeTargets {
		_ := processScrapeTarget(promAPI, state, target)
	}
	return nil
}

func processScrapeTarget(promAPI prometheusv1.API, state *cloudburst.State, scrapeTarget cloudburst.ScrapeTarget) error {
	value, _, err := promAPI.Query(context.TODO(), scrapeTarget.Query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to run processScrapeTargets: %w", err)
	}
	vec := value.(model.Vector)
	var queryResult float64

	for _, v := range vec {
		queryResult = float64(v.Value)
	}

	return scale(queryResult, state, scrapeTarget)
}

func scale(queryResult float64, state *cloudburst.State, scrapeTarget cloudburst.ScrapeTarget) error {

	instances := scrapeTarget.Instances

	sumTerminatingInstances := countTerminatingInstances(instances)
	sumProgressInstances := countInstancesByStatus(instances, cloudburst.Progress)

	var demand = (queryResult + sumTerminatingInstances) - sumProgressInstances

	if demand == 0 {
		return removeAllPending(state, instances)
	}
	if demand > 0 {
		return scaleUp(demand, state, instances)
	}
	if demand < 0 {
		return scaleDown(demand, state, instances)
	}
	return nil
}

func removeAllPending(state *cloudburst.State, instances []cloudburst.Instance) error {
	pendingInstances := getInstancesByStatus(instances, cloudburst.Pending)
	err := state.RemoveInstances(pendingInstances)
	if err != nil {
		return err
	}
}

func scaleUp(demand float64, state *cloudburst.State, instances []cloudburst.Instance) error {
	for i := 0.00; i <= demand; i+1 {

	}
	state.CreateInstances(pending)
}

func scaleDown(demand float64, state *cloudburst.State, instances []cloudburst.Instance) error {

}

func countTerminatingInstances(instances []cloudburst.Instance) float64 {
	var sum int64
	for _, instance := range instances {
		if instance.Active {
			sum++
		}
	}
	return float64(sum)
}

func countInstancesByStatus(instances []cloudburst.Instance, status cloudburst.Status) float64 {
	var sum int64
	for _, instance := range instances {
		if isMatchingStatus(instance, status) {
			sum++
		}
	}
	return float64(sum)
}

func getInstancesByStatus(instances []cloudburst.Instance, status cloudburst.Status) []cloudburst.Instance {
	var res []cloudburst.Instance
	for _, instance := range instances {
		if isMatchingStatus(instance, status) {
			res = append(res, instance)
		}
	}
	return res
}

func isMatchingStatus(instance cloudburst.Instance, status cloudburst.Status) bool {
	return instance.Status.Status == status
}
