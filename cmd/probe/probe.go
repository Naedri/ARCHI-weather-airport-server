package main

import (
	"encoding/json"
	"fmt"
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
	probeType string
	lastRead  time.Time
	id        string
	delta     float64
}

func (probe *Probe) readProbe() (value float64) {
	p := perlin.NewPerlinRandSource(2.0, 2.0, 4, rand.NewSource(int64(3)))
	v := 20 + p.Noise1D(probe.delta)*20
	probe.delta += 0.1
	probe.lastRead = time.Now()
	return v
}

func init() {
	// Register the probe to redis
	utils.HSET("probes", probeID, []byte(probeID))
}

func main() {
	probe := Probe{
		probeType: probeDataType,
		lastRead:  time.Now(),
		id:        probeID,
		delta:     0,
	}
	m := utils.MqttConnection{
		Client: utils.GetDefaultClient(messagePubHandler, connectHandler, connectionLostHandler),
		Topic:  fmt.Sprintf("%s/probe/%s/%s", IATA, probeDataType, probeID),
	}

	for {
		t := time.Now()
		value := model.ProbeMessage{
			Data:     probe.readProbe(),
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
