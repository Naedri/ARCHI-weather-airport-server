package main

import (
	"encoding/json"
	"fmt"
	"meteo_des_aeroports/internal/model"
	"meteo_des_aeroports/internal/utils"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
)

var iataRegistered = false

var (
	qos, _        = strconv.Atoi(os.Getenv("MQTT_QOS"))
	IATA          = os.Getenv("IATA")
	probeDataType = utils.GetDataTypeFromEnv()
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
	toJson := model.ProbeMessage{}
	err := json.Unmarshal([]byte(msg.Payload()), &toJson)
	if err != nil {
		fmt.Println(err.Error())
	}
	//Redis key: ${IATA}:probe:${dataType}:${probeID}
	redisKey := fmt.Sprintf("%s:probe:%s:%s", toJson.IATA, toJson.DataType, toJson.Id)
	fmt.Println(redisKey)
	value := fmt.Sprintf("%.2f", toJson.Data)
	t := toJson.Timestamp
	dateValue, _ := time.Parse("2006-01-02-15-04-05", t)
	dateToUnixMilli := strconv.Itoa(int(dateValue.Unix()))
	utils.ZSet(redisKey, dateToUnixMilli, value)
	if !iataRegistered {
		err := utils.SetAdd(utils.IataListName, []byte(toJson.IATA))
		if err == nil {
			iataRegistered = true
		}
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	client := utils.GetDefaultClient(messagePubHandler, connectHandler, connectionLostHandler)

	topic := fmt.Sprintf("%s/+/%s/%s", IATA, probeDataType, probeID)

	subToken := client.Subscribe(topic, byte(qos), probeDataHandler)

	subToken.Wait()

	fmt.Printf("Subscribed to topic %s", topic)
	<-c

}
