package main

import (
	"flag"
	"fmt"
	"github.com/cbrgm/cloudburst/example/sorting"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"strconv"
	"time"
)


// round((avg_antwortzeit÷slo_antwortzeit)−1)×instanzen_lokal)
// RED Metrics implementation
var (
	requestsDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "example_sorting_request_duration_seconds",
		Help: "The duration of the requests to the sorting service.",
	})

	requestsCurrent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "example_sorting_request_current",
		Help: "The current number of requests to the sorting service.",
	})

	requestsStatus = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "example_sorting_requests_total",
		Help: "The total number of requests to the sorting service by status.",
	}, []string{"status"})

	requestsErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "example_sorting_errors",
		Help: "The total number of errors",
	})
)

func main() {
	var port = flag.String("port", "8080", "help message for flag n")
	flag.Parse()

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	svc := sorting.NewService()
	lsvc := sorting.NewLoggingService(logger, svc)

	api := sorting.NewSortingApi(lsvc)

	r := chi.NewRouter()
	r.Post("/bubblesort", HandleFunc(api.BubbleSortHandler))
	r.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf(":%s", *port)
	level.Info(logger).Log("msg", "webservice is running", "url", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		level.Error(logger).Log("err", err)
	}
}

// HandlerFunc is a wrapper a http handler func.
// Simplifies error handling and metrics calculation for incoming requests.
type HandlerFunc func(http.ResponseWriter, *http.Request) (int, error)

func HandleFunc(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		requestsCurrent.Inc()

		defer func() {
			requestsDuration.Observe(time.Since(now).Seconds())
			requestsCurrent.Dec()
		}()

		statusCode, err := h(w, r)
		requestsStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()

		if err != nil {
			http.Error(w, err.Error(), statusCode)
			requestsErrors.Inc()
			return
		}
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
		}
	}
}