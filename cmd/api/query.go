package main

import (
	"context"
	"fmt"
	"github.com/cbrgm/cloudburst/cloudburst"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"math/rand"
	"strconv"
	"time"
)

func processScrapeTargets(promAPI prometheusv1.API, state *cloudburst.State) error {
	scrapeTargets, err := state.ListScrapeTargets()
	if err != nil {
		return err
	}

	for _, target := range scrapeTargets {
		err = processScrapeTarget(promAPI, state, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func processScrapeTarget(promAPI prometheusv1.API, state *cloudburst.State, scrapeTarget cloudburst.ScrapeTarget) error {
	value, _, err := promAPI.Query(context.TODO(), scrapeTarget.Query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to run processScrapeTargets: %w", err)
	}
	vec := value.(model.Vector)
	var queryResult int

	for _, v := range vec {
		if v.Value.String() == "NaN" {
			queryResult = 0
		} else {
			queryResult = int(v.Value)
		}
	}

	println(queryResult)

	return scale(queryResult, state, scrapeTarget)
}

func scale(queryResult int, state *cloudburst.State, scrapeTarget cloudburst.ScrapeTarget) error {

	instances := scrapeTarget.Instances

	sumTerminatingInstances := countTerminatingInstances(instances)
	sumProgressInstances := countInstancesByStatus(instances, cloudburst.Progress)

	var demand = (queryResult + sumTerminatingInstances) - sumProgressInstances

	if demand == 0 {
		return removeAllPending(state, scrapeTarget)
	}
	if demand > 0 {
		return scaleUp(demand, state, scrapeTarget)
	}
	if demand < 0 {
		return scaleDown(demand, state, scrapeTarget)
	}
	return nil
}

func removeAllPending(state *cloudburst.State, scrapeTarget cloudburst.ScrapeTarget) error {
	pendingInstances := getInstancesByStatus(scrapeTarget.Instances, cloudburst.Pending)
	err := state.RemoveInstances(scrapeTarget.Name, pendingInstances)
	if err != nil {
		return err
	}
	return nil
}

func scaleUp(demand int, state *cloudburst.State, scrapeTarget cloudburst.ScrapeTarget) error {

	pendingInstances := countInstancesByStatus(scrapeTarget.Instances, cloudburst.Pending)

	var pending []cloudburst.Instance
	for i := pendingInstances; i < demand; i++ {
		pending = append(pending, cloudburst.Instance{
			Name:     strconv.Itoa(rand.Intn(100000)),
			Endpoint: "",
			Active:   true,
			Status: cloudburst.InstanceStatus{
				Status:  cloudburst.Pending,
				Started: time.Now(),
			},
		})
	}
	_, err := state.CreateInstances(scrapeTarget.Name, pending)
	return err
}

func scaleDown(demand int, state *cloudburst.State, scrapeTarget cloudburst.ScrapeTarget) error {
	err := removeAllPending(state, scrapeTarget)
	if err != nil {
		return err
	}
	// TODO: label instances to be active = false
	return nil
}

func countTerminatingInstances(instances []cloudburst.Instance) int {
	var sum int
	for _, instance := range instances {
		if instance.Active == false {
			sum++
		}
	}
	return sum
}

func countInstancesByStatus(instances []cloudburst.Instance, status cloudburst.Status) int {
	var sum int
	for _, instance := range instances {
		if isMatchingStatus(instance, status) {
			sum++
		}
	}
	return sum
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
