package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"github.com/Systemnick/argo_exporter/internal/indications/repository/indications"
	"github.com/Systemnick/argo_exporter/pkg/config"
	"github.com/Systemnick/argo_exporter/pkg/prometheus/exporter"
)

func Start(log *zerolog.Logger, cfg *config.Config, indications *indications.MURMetrics) {
	log.Info().Msgf("Starting Argo exporter (Version: %s)", cfg.Version)

	exp, err := exporter.NewArgoExporter(indications)
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(exp)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Argo Exporter (Version ` + cfg.Version + `)</title></head>
			<body>
			<h1>Argo Exporter</h1>
			<p><a href="` + cfg.MetricsPath + `">Metrics</a></p>
			<h2>More information:</h2>
			<p><a href="https://github.com/Systemnick/argo_exporter">github.com/Systemnick/argo_exporter</a></p>
			</body>
			</html>`))
	})
	http.Handle(cfg.MetricsPath, promhttp.Handler())

	log.Info().Msgf("Listening for %s on %s", cfg.MetricsPath, cfg.ListenAddress)
	log.Fatal().Err(http.ListenAndServe(cfg.ListenAddress, nil))
}
