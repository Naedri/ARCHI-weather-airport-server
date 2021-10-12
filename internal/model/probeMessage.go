package model

import "meteo_des_aeroports/internal/utils"

// ${IATA}:probe:${probtype}:${probeId}
type ProbeMessage struct {
	Data      float64
	DataType  utils.DataType
	Timestamp string
	Id        string
	IATA      string
}
