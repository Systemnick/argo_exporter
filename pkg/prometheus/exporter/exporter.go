package exporter

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Systemnick/argo_exporter/internal/indications/repository/indications"
)

type ArgoExporter struct {
	indications *indications.MURMetrics
}

func NewArgoExporter(indications *indications.MURMetrics) (*ArgoExporter, error) {
	return &ArgoExporter{
		indications: indications,
	}, nil
}

func (e *ArgoExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range meterMetrics {
		ch <- metric.desc
	}
}

func (e *ArgoExporter) Collect(ch chan<- prometheus.Metric) {
	for regAddress, regStatus := range e.indications.RegStatus {
		ch <- prometheus.MustNewConstMetric(
			registrarUpDesc,
			prometheus.GaugeValue,
			float64(regStatus),
			regAddress)
	}

	for regAddress, regLastPollStartTime := range e.indications.RegLastPollStartTime {
		ch <- prometheus.MustNewConstMetric(
			registrarLastPollStartTimeDesc,
			prometheus.GaugeValue,
			float64(regLastPollStartTime.Unix()),
			regAddress)
	}

	for regAddress, regLastPollEndTime := range e.indications.RegLastPollEndTime {
		ch <- prometheus.MustNewConstMetric(
			registrarLastPollEndTimeDesc,
			prometheus.GaugeValue,
			float64(regLastPollEndTime.Unix()),
			regAddress)
	}

	for regAddress, duration := range e.indications.RegDuration {
		ch <- prometheus.MustNewConstMetric(
			registrarLastPollDurationDesc,
			prometheus.GaugeValue,
			float64(duration),
			regAddress)
	}

	for meterNumber, meter := range e.indications.Meters {
		v, err := strconv.ParseFloat(meter["PollDuration"], 64)
		if err != nil {
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			meterPollDurationDesc,
			prometheus.GaugeValue,
			v,
			meter["RegistrarIP"], meterNumber, meter["AdapterID"])

		v, err = strconv.ParseFloat(meter["LastIndicationsTime"], 64)
		if err != nil {
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			meterLastIndicationsTimeDesc,
			prometheus.GaugeValue,
			v,
			meter["RegistrarIP"], meterNumber, meter["AdapterID"])
	}

	for _, metric := range meterMetrics {
		for meterNumber, meter := range e.indications.Meters {
			if value, ok := meter[metric.fieldName]; ok {
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					continue
				}

				ch <- prometheus.MustNewConstMetric(
					metric.desc,
					metric.valueType,
					v,
					meterNumber)
			}
		}
	}
}
