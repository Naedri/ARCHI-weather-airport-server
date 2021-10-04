package main

import (
	"encoding/json"
	"fmt"
	"log"
	"meteo_des_aeroports/internal/utils"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
)

type ProbeMessage struct {
	key       string
	data      float64
	dataType  string
	timestamp time.Time
	id        string
}

var qos, _ = strconv.Atoi(os.Getenv("MQTT_QOS"))
var clientID = os.Getenv("MQTT_CLIENT_ID")
var IATA = os.Getenv("IATA")
var probeDataType = os.Getenv("PROBE_DATATYPE")
var probeID = os.Getenv("PROBE_ID")
var conn redis.Conn

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

var probeDataHandler = func(clien mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	toJson := ProbeMessage{}
	err := json.Unmarshal([]byte(msg.Payload()), &toJson)
	fmt.Printf("Operation: %s", toJson.Operation)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = conn.Do("HSET", toJson.key, toJson.data)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	var err error
	conn, err = redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
}

func main() {
	client := utils.SetupClient(clientID, messagePubHandler, connectHandler, connectionLostHandler)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		// TODO change panic and retry to connect
		panic(token.Error())
	}

	topic := fmt.Sprintf("%s/+/%s/%s", IATA, probeDataType, probeID)

	subToken := client.Subscribe(topic, byte(qos), probeDataHandler)

	subToken.Wait()

	fmt.Printf("Subscribed to topic %s", topic)

	for {
	}

}
