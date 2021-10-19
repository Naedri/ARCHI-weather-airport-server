package utils

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
)

type MqttConnection struct {
	Client mqtt.Client
	Topic  string //a client can only publish to an individual topic
}

var brokerURL = os.Getenv("MQTT_BROKER_URL")
var clientID = os.Getenv("MQTT_CLIENT_ID")

//creation and connection of a mqtt Client
func createClient(options *mqtt.ClientOptions) mqtt.Client {
	Client := mqtt.NewClient(options)
	token := Client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Printf("token error : %s\n", token.Error())
		time.Sleep(time.Second * 10)
		token = Client.Connect()
	}
	return Client
}

/*
setting the options for the Client
used to set broker, port, client id, callback and messagePubHandler.
*/
func SetUpClient(brokerURL string, clientID string, pubHand mqtt.MessageHandler, connectHand mqtt.OnConnectHandler, lostHand mqtt.ConnectionLostHandler) *mqtt.ClientOptions {
	options := mqtt.NewClientOptions()
	options.AddBroker(brokerURL)

	options.SetClientID(clientID) // identify each of the clients connecting to the MQTT broker

	options.SetDefaultPublishHandler(pubHand) //  global MQTT pub message processing
	options.OnConnect = connectHand           // callback for the connection
	options.OnConnectionLost = lostHand       //  callback for connection loss

	return options
}

func GetClient(brokerURL string, clientID string, defaultMessagePubHandler mqtt.MessageHandler, connectHandler mqtt.OnConnectHandler, connectionLostHandler mqtt.ConnectionLostHandler) mqtt.Client {
	options := SetUpClient(brokerURL, clientID, defaultMessagePubHandler, connectHandler, connectionLostHandler)
	return createClient(options)
}

// brokerURL and clientID are from the .env file
func GetDefaultClient(defaultMessagePubHandler mqtt.MessageHandler, connectHandler mqtt.OnConnectHandler, connectionLostHandler mqtt.ConnectionLostHandler) mqtt.Client {
	return GetClient(brokerURL, clientID, defaultMessagePubHandler, connectHandler, connectionLostHandler)
}
