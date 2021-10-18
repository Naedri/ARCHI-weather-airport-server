package utils

import (
	"os"
)

type DataType string

const (
	Temperature         DataType = "temperature"
	AtmosphericPressure DataType = "atmospheric"
	WindSpeed           DataType = "wind_speed"
)

var DataTypes = []DataType{
	Temperature,
	AtmosphericPressure,
	WindSpeed,
}

func readDataTypeFromEnv() string {
	return os.Getenv("PROBE_DATATYPE")
}

func convertDataType(dataType string) DataType {
	return DataType(dataType)
}

func GetDataTypeFromEnv() DataType {
	return convertDataType(readDataTypeFromEnv())

}
