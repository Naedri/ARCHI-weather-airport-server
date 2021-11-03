package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"meteo_des_aeroports/internal/model"
	"meteo_des_aeroports/internal/utils"
	"os"
	"strconv"
	"time"

	perlin "github.com/aquilax/go-perlin"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
)

/*
The publisher sends a message on a particular topic.
A publisher is not limited to a single topic, he can send several separate messages on different topics.
*/

var qos, _ = strconv.Atoi(os.Getenv("MQTT_QOS"))
var IATA = os.Getenv("IATA")
var probeDataType = utils.GetDataTypeFromEnv()
var probeID = os.Getenv("PROBE_ID")

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

type Probe struct {
	probeType utils.DataType
	lastRead  time.Time
	id        string
	delta     float64
}

func (probe *Probe) readProbe(probeDataType utils.DataType) (value float64) {
	switch probeDataType {
	case utils.AtmosphericPressure:
		return probe.generateProbeData(1013, 10, 0, 0.1)
	case utils.WindSpeed:
		return probe.generateProbeData(10, 30, 0, 0.1)
	default:
		return probe.generateProbeData(20, 20, math.Inf(-1), 0.001)
	}
}

func (probe *Probe) generateProbeData(average float64, delta float64, min float64, deltaTime float64) (value float64) {
	p := perlin.NewPerlinRandSource(2.0, 2.0, 4, rand.NewSource(int64(3)))
	v := average + p.Noise1D(probe.delta)*delta
	probe.delta += deltaTime
	probe.lastRead = time.Now().UTC()
	return math.Max(min, v)
}

func main() {
	probe := Probe{
		probeType: probeDataType,
		lastRead:  time.Now().UTC(),
		id:        probeID,
		delta:     0,
	}
	m := utils.MqttConnection{
		Client: utils.GetDefaultClient(messagePubHandler, connectHandler, connectionLostHandler),
		Topic:  fmt.Sprintf("%s/probe/%s/%s", IATA, probeDataType, probeID),
	}

	for {
		t := time.Now().UTC()
		value := model.ProbeMessage{
			Data:     probe.readProbe(probeDataType),
			DataType: probeDataType,
			// YYYY-MM-DD-hh-mm-ss
			Timestamp: fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d",
				t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second()),
			Id:   probeID,
			IATA: IATA,
		}

		// ${IATA}:probe:${probtype}:${probeId}
		valueJSONFormated, _ := json.Marshal(value)

		fmt.Printf("%s\n", valueJSONFormated)

		m.Client.Publish(m.Topic, byte(qos), false, valueJSONFormated)

		time.Sleep(time.Second)
		// time.Sleep(time.Second * time.Duration(deltaTime))
	}
}
