package utils

import (
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
)

var brokerURL = os.Getenv("MQTT_BROKER_URL")

func SetupClient(clientID string, defaultMessagePubHandler mqtt.MessageHandler, connectHandler mqtt.OnConnectHandler, connectionLostHandler mqtt.ConnectionLostHandler) mqtt.Client {

	options := mqtt.NewClientOptions()
	options.AddBroker(brokerURL)

	options.SetClientID(clientID)

	options.SetDefaultPublishHandler(defaultMessagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler

	client := mqtt.NewClient(options)

	return client

}
