package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"os"
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

	r := chi.NewRouter()
	r.Post("/bubblesort", HandleFunc(BubbleSortHandler))

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
		statusCode, err := h(w, r)

		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
		}
	}
}

type SortingRequest struct {
	Numbers []int `json:"numbers"`
	Sorted  []int `json:"sorted"`
}

func BubbleSortHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, errors.New("method not allowed")
	}

	resp := SortingRequest{
		Numbers: []int{1, 2, 3, 4},
		Sorted:  []int{1, 2, 3, 4},
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to write http response")
	}
	return http.StatusOK, nil
}
