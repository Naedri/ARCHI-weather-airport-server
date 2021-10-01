package main

import (
	"fmt"
	"math/rand"
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

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

func main() {

	client := utils.SetupClient(clientID, messagePubHandler, connectHandler, connectionLostHandler)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		// TODO change panic and retry to connect
		panic(token.Error())
	}

	topic := fmt.Sprintf("%s/%s/%s", IATA, probeDataType, probeID)

	p := perlin.NewPerlinRandSource(2.0, 2.0, 4, rand.NewSource(int64(3)))
	x := 0.0
	for {

		value := 20 + p.Noise1D(x)*20

		fmt.Printf("%f\n", value)

		client.Publish(topic, byte(qos), false, fmt.Sprintf("%f", value))

		x += 0.01
		time.Sleep(time.Second)
	}
}
