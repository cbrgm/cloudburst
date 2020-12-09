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

func processScrapeTargets(promAPI prometheusv1.API, db State) error {
	scrapeTargets, err := db.ListScrapeTargets()
	if err != nil {
		return err
	}

	for _, target := range scrapeTargets {
		err = processScrapeTarget(promAPI, db, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func processScrapeTarget(promAPI prometheusv1.API, state State, scrapeTarget cloudburst.ScrapeTarget) error {
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

func scale(queryResult int, state State, scrapeTarget cloudburst.ScrapeTarget) error {

	instances, err := state.GetInstancesForTarget(scrapeTarget.Name)
	if err != nil {
		return err
	}

	sumTerminatingInstances := countTerminatingInstances(instances)
	sumProgressInstances := countInstancesByStatus(instances, cloudburst.Progress)

	var demand = (queryResult + sumTerminatingInstances) - sumProgressInstances

	if demand == 0 {
		return removeAllPending(state, instances)
	}
	if demand > 0 {
		return scaleUp(state, demand, scrapeTarget, instances)
	}
	if demand < 0 {
		return scaleDown(state, demand, scrapeTarget, instances)
	}
	return nil
}

func removeAllPending(state State, instances []cloudburst.Instance) error {
	pendingInstances := getInstancesByStatus(instances, cloudburst.Pending)
	err := state.RemoveInstances(pendingInstances)
	if err != nil {
		return err
	}
	return nil
}

func scaleUp(state State, demand int, scrapeTarget cloudburst.ScrapeTarget, instances []cloudburst.Instance) error {

	pendingInstances := countInstancesByStatus(instances, cloudburst.Pending)

	var pending []cloudburst.Instance
	for i := pendingInstances; i < demand; i++ {
		pending = append(pending, cloudburst.Instance{
			Name:     scrapeTarget.Name + "-" + strconv.Itoa(rand.Intn(100000)),
			Target:   scrapeTarget.Name,
			Endpoint: "",
			Active:   true,
			Status: cloudburst.InstanceStatus{
				Status:  cloudburst.Pending,
				Started: time.Now(),
			},
		})
	}
	_, err := state.SaveInstances(scrapeTarget.Name, pending)
	return err
}

func scaleDown(state State, demand int, scrapeTarget cloudburst.ScrapeTarget, instances []cloudburst.Instance) error {
	err := removeAllPending(state, instances)
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
