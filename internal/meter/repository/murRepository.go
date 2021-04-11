package repository

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/Systemnick/argo_exporter/internal/domain/meter"
	"github.com/Systemnick/argo_exporter/internal/meter/useCase"
)

type MURRepo struct{}

func (MURRepo) GetCurrentIndication(registrarId uuid.UUID, meterId int16) ([]useCase.MeterInfo, error) {
	panic("implement me")
}

func parseFlatIndication(indications map[string]string) []useCase.MeterInfo {
	heatMeterInfo := useCase.MeterInfo{
		Meter: meter.Meter{
			Purpose: meter.Heat,
		},
	}
	hotWaterMeterInfo := useCase.MeterInfo{
		Meter: meter.Meter{
			Purpose: meter.WaterHot,
		},
	}
	coldWaterMeterInfo := useCase.MeterInfo{
		Meter: meter.Meter{
			Purpose: meter.WaterCold,
		},
	}
	hotWaterMeter2Info := useCase.MeterInfo{
		Meter: meter.Meter{
			Purpose: meter.WaterHot,
		},
	}
	coldWaterMeter2Info := useCase.MeterInfo{
		Meter: meter.Meter{
			Purpose: meter.WaterCold,
		},
	}

	meterInfos := []useCase.MeterInfo{}

	for k, v := range indications {
		var currentMeter *useCase.MeterInfo

		switch k {
		case "Adapter":
			heatMeterInfo.Meter.Number = meter.Number(v)
		case "Energy":
			currentMeter = &heatMeterInfo
		case "Volume1":
			currentMeter = &hotWaterMeterInfo
		case "Volume2":
			currentMeter = &coldWaterMeterInfo
		case "Volume3":
			currentMeter = &hotWaterMeter2Info
		case "Volume4":
			currentMeter = &coldWaterMeter2Info
		}

		if currentMeter == nil {
			continue
		}

		value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			fmt.Printf("bad integer: %s\n", v)
		}

		currentMeter.CurrentValue = value
		meterInfos = append(meterInfos, *currentMeter)
	}

	return meterInfos
}
