package main

import prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"

type AutoScaler struct {
	prom  prometheusv1.API
	state State
}

func NewAutoScaler() {

}
