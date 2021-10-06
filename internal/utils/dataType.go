package utils

import (
	"os"
	"strconv"
)

type DataType int

const (
	Temperature DataType = iota
	AtmosphericPressure
	WindSpeed
)

func readDataTypeFromEnv() int {
	i, err := strconv.Atoi(os.Getenv("PROBE_DATATYPE"))
	if err != nil || i < 0 {
		i = 0
	}
	return i
}

func GetDataTypeFromEnv() string {
	i := readDataTypeFromEnv()
	return DataType(i).String()
}

func (d DataType) String() string {
	return [...]string{"temperature", "atmospheric_pressure", "wind_speed"}[d]
}
