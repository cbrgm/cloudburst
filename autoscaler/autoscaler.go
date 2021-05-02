package autoscaler

import (
	"context"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/cbrgm/cloudburst/metrics"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"time"
)

type State interface {
	ListScrapeTargets() ([]*cloudburst.ScrapeTarget, error)
	GetInstance(scrapeTarget string, name string) (*cloudburst.Instance, error)
	GetInstances(scrapeTarget string) ([]*cloudburst.Instance, error)
	RemoveInstances(scrapeTarget string, instances []*cloudburst.Instance) error
	RemoveInstance(scrapeTarget string, instance *cloudburst.Instance) error
	SaveInstances(scrapeTarget string, instances []*cloudburst.Instance) ([]*cloudburst.Instance, error)
	SaveInstance(scrapeTarget string, instance *cloudburst.Instance) (*cloudburst.Instance, error)
}

type ScalingMetrics interface {
	SetReceiverMetricValue(service string, value float64)
	SetCalculatorInstancesTotal(service, provider, status string, value float64)
	SetCalculatorInstanceDemandResult(service, provider string, value float64)
}

type ScalingCalculator interface {
	Calculate(scrapeTarget *cloudburst.ScrapeTarget, instances []*cloudburst.Instance, metricValue float64) cloudburst.ScalingResult
}

type ScalingReceiver interface {
	Poll(query string) (float64, error)
	PollFrom(url string, query string) (float64, error)
}

type ScalingPreparator interface {
	Prepare(result cloudburst.ScalingResult, scrapeTarget *cloudburst.ScrapeTarget) error
}

type Scaling struct {
	logger    log.Logger
	startTime time.Time

	state     State
	metrics   ScalingMetrics
	receive   ScalingReceiver
	calculate ScalingCalculator
	prepare   ScalingPreparator

	interval time.Duration
	url      string
}

// ScalingOption passed to NewBot to change the default instance.
type ScalingOption func(b *Scaling) error

func NewScaling(state State, url string, opts ...ScalingOption) (*Scaling, error) {
	sMetrics := metrics.NewDefaultPrometheus()
	receive, err := cloudburst.NewMetricsReceiver(url)
	if err != nil {
		return nil, err
	}
	calculate := cloudburst.NewReactiveThreshold()
	prepare := newPreparator(state)

	return NewScalingWithOptions(state, sMetrics, receive, calculate, prepare, url, opts...)
}

func NewScalingWithOptions(state State, sMetrics ScalingMetrics, receive ScalingReceiver, calculate ScalingCalculator, prepare ScalingPreparator, url string, opts ...ScalingOption) (*Scaling, error) {
	s := &Scaling{
		logger:    log.NewNopLogger(),
		startTime: time.Now(),
		state:     state,
		metrics:   sMetrics,
		receive:   receive,
		calculate: calculate,
		prepare:   prepare,
		interval:  30 * time.Second,
		url:       url,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func WithLogger(l log.Logger) ScalingOption {
	return func(s *Scaling) error {
		s.logger = l
		return nil
	}
}

func WithMetrics(m ScalingMetrics) ScalingOption {
	return func(s *Scaling) error {
		s.metrics = m
		return nil
	}
}

func WithReceiver(r ScalingReceiver) ScalingOption {
	return func(s *Scaling) error {
		s.receive = r
		return nil
	}
}

func WithCalculator(c ScalingCalculator) ScalingOption {
	return func(s *Scaling) error {
		s.calculate = c
		return nil
	}
}

func WithInterval(t time.Duration) ScalingOption {
	return func(s *Scaling) error {
		s.interval = t
		return nil
	}
}

func (s *Scaling) Run(ctx context.Context) error {
	tick := time.NewTicker(s.interval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			err := s.processScrapeTargets()
			if err != nil {
				level.Info(s.logger).Log(
					"msg", "failed to process scrapetargets",
					"error", err,
				)
			}
		}
	}
}

func (s *Scaling) processScrapeTargets() error {
	scrapeTargets, err := s.state.ListScrapeTargets()
	if err != nil {
		return err
	}
	for _, target := range scrapeTargets {
		err = s.processScrapeTarget(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scaling) processScrapeTarget(scrapeTarget *cloudburst.ScrapeTarget) error {
	// todo: handle instances with state failure
	err := s.cleanupTerminatedInstances(scrapeTarget)
	if err != nil {
		return err
	}

	// receive metric values from prometheus
	metricValue, err := s.receive.PollFrom(s.url, scrapeTarget.Query)
	if err != nil {
		return err
	}
	s.metrics.SetReceiverMetricValue(scrapeTarget.Name, metricValue)
	return s.ProcessScrapeTargetWithValue(scrapeTarget, metricValue)
}

func (s *Scaling) ProcessScrapeTargetWithValue(scrapeTarget *cloudburst.ScrapeTarget, value float64) error {
	// get all instances from state, only call once!
	instances, err := s.state.GetInstances(scrapeTarget.Name)
	if err != nil {
		return err
	}
	s.updateInstancesTotalMetric(scrapeTarget, instances)

	// calculate demand
	result := s.calculate.Calculate(scrapeTarget, instances, value)
	s.updateInstanceDemandMetric(scrapeTarget, result)

	// prepare execution
	err = s.prepare.Prepare(result, scrapeTarget)
	if err != nil {
		return err
	}
	return nil
}

func (s *Scaling) cleanupTerminatedInstances(target *cloudburst.ScrapeTarget) error {
	instances, err := s.state.GetInstances(target.Name)
	if err != nil {
		return err
	}

	for _, instance := range instances {
		if instance.Status.Status == cloudburst.Terminated || instance.Status.Status == cloudburst.Failure {
			err := s.state.RemoveInstance(target.Name, instance)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Scaling) updateInstanceDemandMetric(scrapeTarget *cloudburst.ScrapeTarget, result cloudburst.ScalingResult) {
	for _, result := range result.Result {
		s.metrics.SetCalculatorInstanceDemandResult(scrapeTarget.Name, result.Provider, float64(result.InstanceDemand))
	}
}

func (s *Scaling) updateInstancesTotalMetric(scrapeTarget *cloudburst.ScrapeTarget, instances []*cloudburst.Instance) {
	for provider, _ := range scrapeTarget.ProviderSpec.Weights {
		in := cloudburst.GetInstancesByProvider(instances, provider)
		count := make(map[cloudburst.Status]float64)
		for _, instance := range in {
			if val, ok := count[instance.Status.Status]; ok {
				count[instance.Status.Status] = val + 1
			} else {
				count[instance.Status.Status] = 1
			}
		}
		s.metrics.SetCalculatorInstancesTotal(scrapeTarget.Name, provider, string(cloudburst.Pending), count[cloudburst.Pending])
		s.metrics.SetCalculatorInstancesTotal(scrapeTarget.Name, provider, string(cloudburst.Running), count[cloudburst.Running])
		s.metrics.SetCalculatorInstancesTotal(scrapeTarget.Name, provider, string(cloudburst.Terminated), count[cloudburst.Terminated])
		s.metrics.SetCalculatorInstancesTotal(scrapeTarget.Name, provider, string(cloudburst.Failure), count[cloudburst.Failure])
	}
}
