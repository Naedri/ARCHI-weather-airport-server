package utils

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
)

type MqttConnection struct {
	client mqtt.Client
	topic  string
}

var brokerURL = os.Getenv("MQTT_BROKER_URL")
var clientID = os.Getenv("MQTT_BROKER_URL")

//creation and connection of a mqtt client
func createClient(options *mqtt.ClientOptions) mqtt.Client {
	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Printf("token error : %s\n", token.Error())
		time.Sleep(time.Second * 10)
		token = client.Connect()
	}
	return client
}

//setting the options for the client
func setUpClient(brokerURL string, clientID string, pubHand mqtt.MessageHandler, connectHand mqtt.OnConnectHandler, lostHand mqtt.ConnectionLostHandler) *mqtt.ClientOptions {

	options := mqtt.NewClientOptions()
	options.AddBroker(brokerURL)

	options.SetClientID(clientID)

	options.SetDefaultPublishHandler(pubHand)
	options.OnConnect = connectHand
	options.OnConnectionLost = lostHand

	return options
}

func GetClient(brokerURL string, clientID string, defaultMessagePubHandler mqtt.MessageHandler, connectHandler mqtt.OnConnectHandler, connectionLostHandler mqtt.ConnectionLostHandler) mqtt.Client {
	options := setUpClient(brokerURL, clientID, defaultMessagePubHandler, connectHandler, connectionLostHandler)
	return createClient(options)
}

func GetDefaultClient(defaultMessagePubHandler mqtt.MessageHandler, connectHandler mqtt.OnConnectHandler, connectionLostHandler mqtt.ConnectionLostHandler) mqtt.Client {
	return GetClient(brokerURL, clientID, defaultMessagePubHandler, connectHandler, connectionLostHandler)
}
