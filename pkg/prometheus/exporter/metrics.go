package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// TODO: Метрики:
	//  - Сколько попыток опроса счётчиков было суммарно
	//  - Время последнего опроса

	registrarUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName("registrar", "", "up"),
		"Whether scraping registrar's metrics was successful.",
		[]string{"registrar_ip"},
		nil)

	registrarPollDesc = prometheus.NewDesc(
		prometheus.BuildFQName("registrar", "", "last_poll_duration_nanoseconds"),
		"Registrar poll duration in nanoseconds.",
		[]string{"registrar_ip"},
		nil)

	meterPollDesc = prometheus.NewDesc(
		prometheus.BuildFQName("meter", "", "last_poll_duration_nanoseconds"),
		"Meter poll duration in nanoseconds.",
		[]string{"registrar_ip", "meter_serial_number", "meter_number_on_registrar"},
		nil)

	meterMetrics = []*meterMetric{
		newMeterMetric(
			"status",
			"Current meter status.",
			prometheus.GaugeValue,
			[]string{"meter_serial_number"},
			"Status"),
		newMeterMetric(
			"energy_gcal",
			"Total amount of energy spent.",
			prometheus.CounterValue,
			[]string{"meter_serial_number"},
			"Energy"),
		newMeterMetric(
			"volume0_cbm",
			"Total general volume amount in cubic meters.",
			prometheus.CounterValue,
			[]string{"meter_serial_number"},
			"Volume"),
		newMeterMetric(
			"volume1_cbm",
			"Total additional volume 1 amount in cubic meters.",
			prometheus.CounterValue,
			[]string{"meter_serial_number"},
			"Volume1"),
		newMeterMetric(
			"volume2_cbm",
			"Total additional volume 2 amount in cubic meters.",
			prometheus.CounterValue,
			[]string{"meter_serial_number"},
			"Volume2"),
		newMeterMetric(
			"volume3_cbm",
			"Total additional volume 3 amount in cubic meters.",
			prometheus.CounterValue,
			[]string{"meter_serial_number"},
			"Volume3"),
		newMeterMetric(
			"volume4_cbm",
			"Total additional volume 4 amount in cubic meters.",
			prometheus.CounterValue,
			[]string{"meter_serial_number"},
			"Volume4"),
		newMeterMetric(
			"consumption_cbm_h",
			"Actual consumption in cubic meters per hour.",
			prometheus.GaugeValue,
			[]string{"meter_serial_number"},
			"Consumption"),
		newMeterMetric(
			"power_kwatt",
			"Actual power in kilowatts.",
			prometheus.GaugeValue,
			[]string{"meter_serial_number"},
			"Power"),
		newMeterMetric(
			"temperature_in_celsius",
			"Actual incoming flow temperature in Celsius.",
			prometheus.GaugeValue,
			[]string{"meter_serial_number"},
			"TemperatureIn"),
		newMeterMetric(
			"temperature_out_celsius",
			"Actual outgoing flow temperature in Celsius.",
			prometheus.GaugeValue,
			[]string{"meter_serial_number"},
			"TemperatureOut"),
		newMeterMetric(
			"runtime_hour",
			"Total meter runtime in hours.",
			prometheus.CounterValue,
			[]string{"meter_serial_number"},
			"Runtime"),
		newMeterMetric(
			"runtime_error_hour",
			"Total meter runtime with error in hours.",
			prometheus.CounterValue,
			[]string{"meter_serial_number"},
			"RuntimeWithError"),
	}
)

type meterMetric struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	fieldName string
}

func newMeterMetric(name string, description string, valueType prometheus.ValueType, labels []string, fieldName string) *meterMetric {
	return &meterMetric{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName("meter", "", name),
			description,
			labels,
			nil),
		valueType: valueType,
		fieldName: fieldName,
	}
}
