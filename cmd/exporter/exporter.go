package main

import (
	"encoding/json"
	"fmt"
	"log"
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

/*
The subscribers are the clients of the broker.
A subscriber subscribes to one (or more) topic in order to be notified of the arrival of new messages on the said topic (s).
Purpose : to create .csv files stored in the local files of the server
*/
var (
	qos, _ = strconv.Atoi(os.Getenv("MQTT_QOS"))
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
	t := toJson.Timestamp
	dateValue, _ := time.Parse("2006-01-02-15-04-05", t)
	fileName := fmt.Sprintf("%s-%s-%d-%02d-%02d.csv", toJson.DataType, toJson.IATA, dateValue.Year(), dateValue.Month(), dateValue.Day())
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	value := fmt.Sprintf("%s,%0.2f\n", toJson.Id, toJson.Data)
	f.Write([]byte(value))
	defer f.Close()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	client := utils.GetDefaultClient(messagePubHandler, connectHandler, connectionLostHandler)

	topic := "+/probe/+/+"

	subToken := client.Subscribe(topic, byte(qos), probeDataHandler)

	subToken.Wait()

	fmt.Printf("Subscribed to topic %s", topic)
	<-c
}
