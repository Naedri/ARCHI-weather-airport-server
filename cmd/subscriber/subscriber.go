package main

import (
	"encoding/json"
	"fmt"
	"meteo_des_aeroports/internal/utils"
	"os"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
)

type ProbeMessage struct {
	Data      string
	DataType  string
	Timestamp string
	Id        string
	IATA      string
}

var (
	qos, _        = strconv.Atoi(os.Getenv("MQTT_QOS"))
	clientID      = os.Getenv("MQTT_CLIENT_ID")
	IATA          = os.Getenv("IATA")
	probeDataType = os.Getenv("PROBE_DATATYPE")
	probeID       = os.Getenv("PROBE_ID")
)
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

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
	if err != nil {
		fmt.Println(err.Error())
	}
	//Redis key: ${IATA}:probe:${dataType}:${probeID}
	redisKey := fmt.Sprintf("%s:probe:%s:%s", toJson.IATA, toJson.DataType, toJson.Id)
	fmt.Println(redisKey)
	toFloat, _ := strconv.ParseFloat(toJson.Data, 64)
	value := fmt.Sprintf("%.2f:%s", toFloat, toJson.Timestamp)
	utils.ZSet(redisKey, toJson.Timestamp, value)
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
