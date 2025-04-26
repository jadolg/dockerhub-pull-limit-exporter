package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const prefix = "dockerhub_pull_"

var (
	pullLimit = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%slimit", prefix),
			Help: "The rate limit for Docker Hub pulls",
		},
		[]string{"account", "source"},
	)
	pullRemaining = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%sremaining", prefix),
			Help: "The remaining pulls for Docker Hub",
		},
		[]string{"account", "source"},
	)
	limitWindowSeconds = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%slimit_window_seconds", prefix),
			Help: "The time window in seconds to which the limit applies",
		},
		[]string{"account", "source"},
	)
	remainingWindowSeconds = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%sremaining_window_seconds", prefix),
			Help: "The time window in seconds to which the remaining pulls applies",
		},
		[]string{"account", "source"},
	)
	errors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%serrors_total", prefix),
			Help: "Exporter errors",
		},
		[]string{"account"},
	)
)

func healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		log.Errorf("error responding to request %v", err)
	}
}

func startMetricsServer(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", healthcheckHandler)
	log.Printf("Starting metrics server on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
