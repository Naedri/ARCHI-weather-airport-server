package model

// ${IATA}:probe:${probtype}:${probeId}
type ProbeMessage struct {
	Data      float64
	DataType  string
	Timestamp string
	Id        string
	IATA      string
}
