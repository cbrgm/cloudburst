package cloudburst

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"time"
)

type MetricsReceiver struct {
	api prometheusv1.API
}

func NewMetricsReceiver(prometheusUrl string) (*MetricsReceiver, error) {
	client, err := api.NewClient(api.Config{Address: prometheusUrl})
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}
	return &MetricsReceiver{
		api: prometheusv1.NewAPI(client),
	}, nil
}

func (m *MetricsReceiver) Poll(query string) (float64, error) {
	value, _, err := m.api.Query(context.TODO(), query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to query metrics: %w", err)
	}
	vec := value.(model.Vector)
	var metricValue float64

	for _, v := range vec {
		if v.Value.String() == "NaN" {
			metricValue = 0
		} else {
			metricValue = float64(v.Value)
		}
	}
	return metricValue, nil
}

func (m *MetricsReceiver) PollFrom(url string, query string) (float64, error) {
	client, err := api.NewClient(api.Config{Address: url})
	if err != nil {
		return 0, fmt.Errorf("failed to create Prometheus client: %w", err)
	}

	promAPI := prometheusv1.NewAPI(client)

	value, _, err := promAPI.Query(context.TODO(), query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to query metrics: %w", err)
	}
	vec := value.(model.Vector)
	var metricValue float64

	for _, v := range vec {
		if v.Value.String() == "NaN" {
			metricValue = 0
		} else {
			metricValue = float64(v.Value)
		}
	}
	return metricValue, nil
}
