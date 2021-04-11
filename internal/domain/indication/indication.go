package indication

import (
	"fmt"
	"reflect"
	"time"
)

type Indication struct {
	AdapterID        uint16
	AdapterNumber    uint32
	AdapterVersion   uint8
	Status           uint8
	CurrentDate      time.Time
	Energy           float64
	Volume           float64
	Volume1          float64
	Volume2          float64
	Volume3          float64
	Volume4          float64
	Consumption      float64
	Power            float64
	TemperatureIn    float64
	TemperatureOut   float64
	Runtime          float64
	RuntimeWithError float64
}

type IndicationMap map[string]string

func (i Indication) String() string {
	s := "{\n"

	v := reflect.ValueOf(i)

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		s += fmt.Sprintf("  %s: %v\n", v.Type().Field(i).Name, v.Field(i).Interface())
	}

	s += "}"

	return s
}
