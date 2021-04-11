package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServiceName    string
	Version        string
	ListenAddress  string
	MetricsPath    string
	ScrapeInterval time.Duration
}

func New(serviceName, version string) (*Config, error) {
	listenAddress := getEnvString("LISTEN_ADDRESS", ":8080")
	metricsPath := getEnvString("METRICS_PATH", "/metrics")
	scrapeInterval := getEnvInt("SCRAPE_INTERVAL", 60) // Minutes

	return &Config{
		ServiceName:    serviceName,
		Version:        version,
		ListenAddress:  listenAddress,
		MetricsPath:    metricsPath,
		ScrapeInterval: time.Duration(scrapeInterval) * time.Minute,
	}, nil
}

func getEnvString(key, defaultValue string) string {
	s := os.Getenv(key)
	if len(s) == 0 {
		return defaultValue
	}

	return s
}

func getEnvInt(key string, defaultValue int) int {
	s := os.Getenv(key)
	if len(s) == 0 {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return i
}
