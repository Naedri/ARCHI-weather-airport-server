package model

import "time"

// ${IATA}:probe:${probtype}:${probeId}
type ProbeMessage struct {
	Key       string
	Data      float64
	DataType  string
	Timestamp time.Time
	Id        string
}
