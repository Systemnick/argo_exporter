package meter

import "time"

type Number string

type Purpose uint8

const (
	WaterCold Purpose = iota
	WaterHot
	Heat
	Electricity
)

type CalibrationDate time.Time

type Meter struct {
	Number              Number
	Purpose             Purpose
	LastCalibrationDate CalibrationDate
	NextCalibrationDate CalibrationDate
}
