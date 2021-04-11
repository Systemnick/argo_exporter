package useCase

import "github.com/Systemnick/argo_exporter/internal/domain/meter"

type MeterInfo struct {
	Meter        meter.Meter
	CurrentValue uint64
}

type FlatIndication struct {
	Indications []MeterInfo
}
