package main

import (
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
var clientID = os.Getenv("MQTT_BROKER_URL")
var IATA = os.Getenv("IATA")
var probeDataType = os.Getenv("PROBE_DATATYPE")
var probeID = os.Getenv("PROBE_ID")
var deltaTime = 10

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

type MqttConnection struct {
	client mqtt.Client
	topic  string
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

func main() {

	var m *MqttConnection = &MqttConnection{}
	probe := Probe{probeType: probeDataType, lastRead: time.Now(), id: probeID, delta: 0}
	m.client = utils.SetupClient(clientID, messagePubHandler, connectHandler, connectionLostHandler)

	token := m.client.Connect()
	for token.Wait() && token.Error() != nil {
		token = m.client.Connect()
	}

	m.topic = fmt.Sprintf("%s/probe/%s/%s", IATA, probeDataType, probeID)

	for {
		value := model.ProbeMessage{
			Key:       m.topic,
			Data:      probe.readProbe(),
			DataType:  probeDataType,
			Timestamp: time.Now(),
		}

		// ${IATA}:probe:${probtype}:${probeId}
		valueJSONFormated := fmt.Sprintf(`
		{
			"Key":"%s",
			"Data":"%f",
			"DataType":"%s",
			"Timestamp":"%s",
			"Id": "%s"
		}`, value.Key, value.Data, value.DataType, value.Timestamp, probeID)

		fmt.Printf("%s\n", valueJSONFormated)

		m.client.Publish(m.topic, byte(qos), false, valueJSONFormated)

		time.Sleep(time.Second)
		// time.Sleep(time.Second * time.Duration(deltaTime))
	}
}
