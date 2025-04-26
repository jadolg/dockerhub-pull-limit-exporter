package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
	"os"
	"time"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	var port int
	var configFile string
	var logLevel string
	var version bool

	flag.IntVar(&port, "port", 9101, "Port to listen on")
	flag.StringVar(&configFile, "config", "config.yaml", "Path to config file")
	flag.StringVar(&logLevel, "loglevel", "info", "Log level")
	flag.BoolVar(&version, "version", false, "prints version and exits")
	flag.Parse()

	err := configureLogs(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"Version": Version,
		"Commit":  Commit,
		"Date":    Date,
	}).Info("Docker Hub Pull Limits Exporter")

	if version {
		return
	}

	config, err := getConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	for _, credential := range config.Credentials {
		log.WithFields(log.Fields{
			"username": credential.Username,
		}).Info("Starting metrics collector")
		ticker := time.NewTicker(config.UpdateInterval)

		go func() {
			for ; true; <-ticker.C {
				log.WithFields(log.Fields{
					"username": credential.Username,
				}).Debug("Collecting metrics")
				err := collectMetrics(credential, config.Timeout)
				if err != nil {
					log.WithFields(log.Fields{
						"username": credential.Username,
					}).Error(err)
					errorsCount.WithLabelValues(credential.Username).Inc()
				} else {
					log.WithFields(log.Fields{
						"username": credential.Username,
					}).Debug("Successfully collected metrics")
				}
			}
		}()
	}

	if err := startMetricsServer(port); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}

func collectMetrics(credential credentials, timeout time.Duration) error {
	token, err := getToken(credential.Username, credential.Password, timeout)
	if err != nil {
		return err
	}
	limit, remaining, limitWindow, remainingWindow, source, err := getLimits(token, timeout)
	if err != nil {
		return err
	}

	pullLimit.WithLabelValues(credential.Username, source).Set(float64(limit))
	pullRemaining.WithLabelValues(credential.Username, source).Set(float64(remaining))
	limitWindowSeconds.WithLabelValues(credential.Username, source).Set(float64(limitWindow))
	remainingWindowSeconds.WithLabelValues(credential.Username, source).Set(float64(remainingWindow))

	return nil
}

func configureLogs(logLevel string) error {
	parsedLogLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	log.SetLevel(parsedLogLevel)
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		log.SetFormatter(&log.JSONFormatter{})
	}
	return err
}
