package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
	"time"
)

const (
	promNamespace           = "cloudburst"
	promApiSubsystem        = "api"
	promReceiverSubsystem   = "receiver"
	promCalculatorSubsystem = "calculator"
	// promPreparatorSubsystem = "preparator"
)

// Prometheus implements the prometheus metrics backend.
type Prometheus struct {
	apiRequestDurationM             *prometheus.HistogramVec
	apiRequestsTotalM               *prometheus.CounterVec
	apiEventDurationM               *prometheus.HistogramVec
	apiEventSubscribersTotal        *prometheus.GaugeVec
	receiverMetricValueM            *prometheus.GaugeVec
	calculatorInstancesTotalM       *prometheus.GaugeVec
	calculatorInstanceDemandResultM *prometheus.GaugeVec

	opts     Options
	registry *prometheus.Registry
	handler  http.Handler
}

func NewDefaultPrometheus() *Prometheus {
	return NewPrometheus(DefaultOptions())
}

// NewPrometheus returns a new Prometheus metric backend.
func NewPrometheus(opts Options) *Prometheus {
	namespace := promNamespace
	if opts.Prefix != "" {
		namespace = strings.TrimSuffix(opts.Prefix, ".")
	}

	// api
	apiRequestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: promApiSubsystem,
		Name:      "http_request_duration_seconds",
		Help:      "http latency to openapi http handlers",
	}, []string{"code", "method", "name"})

	apiRequestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: promApiSubsystem,
		Name:      "http_request_total",
		Help:      "http requests to openapi http handlers",
	}, []string{"code", "method", "name"})

	apiEventDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   namespace,
		Subsystem:   promApiSubsystem,
		Name:        "http_event_duration_seconds",
		Help:        "Duration and error code for server sent events",
		ConstLabels: prometheus.Labels{"events": "instances"},
	}, []string{"status"})

	apiEventSubscribersTotal := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Subsystem:   promApiSubsystem,
		Name:        "http_event_subscriptions_total",
		Help:        "Number of current subscribers",
		ConstLabels: prometheus.Labels{"events": "instances"},
	}, []string{})

	// receiver
	receiverMetricValue := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: promReceiverSubsystem,
		Name:      "metric_value",
		Help:      "metric value received from prometheus for the current scaling iteration",
	}, []string{"target"})

	// calculator
	calculatorInstancesTotal := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: promCalculatorSubsystem,
		Name:      "instances_total",
		Help:      "instances total",
	}, []string{"target", "provider", "status"})

	calculatorInstanceDemandResult := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: promCalculatorSubsystem,
		Name:      "instances_demand",
		Help:      "calculated instance demand for the current scaling iteration",
	}, []string{"target", "provider"})

	p := &Prometheus{
		apiRequestDurationM:             apiRequestDuration,
		apiRequestsTotalM:               apiRequestsTotal,
		apiEventDurationM:               apiEventDuration,
		apiEventSubscribersTotal:        apiEventSubscribersTotal,
		receiverMetricValueM:            receiverMetricValue,
		calculatorInstancesTotalM:       calculatorInstancesTotal,
		calculatorInstanceDemandResultM: calculatorInstanceDemandResult,
		opts:                            opts,
		registry:                        opts.PrometheusRegistry,
		handler:                         nil,
	}

	if p.registry == nil {
		p.registry = prometheus.NewRegistry()
	}
	p.registerMetrics()
	return p
}

func (p *Prometheus) registerMetrics() {
	p.registry.MustRegister(p.apiRequestDurationM)
	p.registry.MustRegister(p.apiRequestsTotalM)
	p.registry.MustRegister(p.apiEventDurationM)
	p.registry.MustRegister(p.apiEventSubscribersTotal)
	p.registry.MustRegister(p.receiverMetricValueM)
	p.registry.MustRegister(p.calculatorInstancesTotalM)
	p.registry.MustRegister(p.calculatorInstanceDemandResultM)

	if p.opts.EnableRuntimeMetrics {
		p.registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
		p.registry.MustRegister(prometheus.NewGoCollector())
	}
}

func (p *Prometheus) CreateHandler() http.Handler {
	return promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{})
}

func (p *Prometheus) getHandler() http.Handler {
	if p.handler != nil {
		return p.handler
	}
	p.handler = p.CreateHandler()
	return p.handler
}

// RegisterHandler satisfies Metrics interface.
func (p *Prometheus) RegisterHandler(path string, mux *http.ServeMux) {
	promHandler := p.getHandler()
	mux.Handle(path, promHandler)
}

// sinceStart returns the seconds passed since the start time until now.
func (p *Prometheus) sinceStart(start time.Time) float64 {
	return time.Since(start).Seconds()
}

func (p *Prometheus) MeasureOpenApiRequestDuration(code, method, name string, start time.Time) {
	t := p.sinceStart(start)
	p.apiRequestDurationM.WithLabelValues(code, method, name).Observe(t)
}
func (p *Prometheus) IncOpenApiRequestsTotal(code, method, name string) {
	p.apiRequestsTotalM.WithLabelValues(code, method, name).Inc()
}

func (p *Prometheus) MeasureApiEventDuration(error string, start time.Time) {
	t := p.sinceStart(start)
	p.apiEventDurationM.WithLabelValues("error").Observe(t)
}

func (p *Prometheus) IncEventSubscribers() {
	p.apiEventSubscribersTotal.WithLabelValues().Inc()
}

func (p *Prometheus) DecEventSubscribers() {
	p.apiEventSubscribersTotal.WithLabelValues().Dec()
}

func (p *Prometheus) SetReceiverMetricValue(service string, value float64) {
	p.receiverMetricValueM.WithLabelValues(service).Set(value)
}

func (p *Prometheus) SetCalculatorInstancesTotal(service, provider, status string, value float64) {
	p.calculatorInstancesTotalM.WithLabelValues(service, provider, status).Set(value)
}

func (p *Prometheus) SetCalculatorInstanceDemandResult(service, provider string, value float64) {
	p.calculatorInstanceDemandResultM.WithLabelValues(service, provider).Set(value)
}
