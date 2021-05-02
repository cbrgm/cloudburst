package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/cbrgm/cloudburst/cloudburst/convert"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"time"
)

type Events interface {
	SubscribeToInstanceEvents(channel chan cloudburst.InstanceEvent) cloudburst.Subscription
	UnsubscribeFromInstanceEvents(s cloudburst.Subscription)
}

func instanceEventsHandler(logger log.Logger, apiMetrics ApiMetrics, events Events) func(w http.ResponseWriter, r *http.Request) {
	observeEvent := func(start time.Time, err error) {
		if err != nil {
			level.Warn(logger).Log(
				"msg", "failed to send server sent event",
				"err", err,
			)
			apiMetrics.MeasureApiEventDuration("error", start)
		} else {
			level.Debug(logger).Log(
				"msg", "successfully sent server sent event",
			)
			apiMetrics.MeasureApiEventDuration("success", start)
		}
	}

	observeSubscription := func() {
		apiMetrics.IncEventSubscribers()
		level.Debug(logger).Log("msg", "subscriber to events")
	}

	observeUnsubscription := func() {
		apiMetrics.DecEventSubscribers()
		level.Debug(logger).Log("msg", "unsubscriber from events")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, `{"err":"server sent events unsupported"}`, http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		instanceEvents := make(chan cloudburst.InstanceEvent, 25)
		subscriptions := events.SubscribeToInstanceEvents(instanceEvents)
		observeSubscription()

		defer func() {
			events.UnsubscribeFromInstanceEvents(subscriptions)
			observeUnsubscription()
		}()

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				close(instanceEvents)
				return
			case event := <-instanceEvents:
				start := time.Now()
				payload, err := json.Marshal(convert.InstanceEventToOpenAPI(event))
				if err != nil {
					observeEvent(start, err)
					continue
				}

				_, err = fmt.Fprintf(w, "data: %s\n\n", payload)
				if err != nil {
					observeEvent(start, err)
					continue
				}
				flusher.Flush()
				observeEvent(start, nil)
			}
		}
	}
}
