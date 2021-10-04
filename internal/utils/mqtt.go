package utils

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttConnection struct {
	client mqtt.Client
	topic  string
}

//creation of a mqtt client
func CreateClient(options *mqtt.ClientOptions) (mqtt.Client, error) {
	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return client, token.Error()
	}
	return client, nil
}

//setting the options for the client
func SetupClient2(clientID string, defaultMessagePubHandler mqtt.MessageHandler, connectHandler mqtt.OnConnectHandler, connectionLostHandler mqtt.ConnectionLostHandler) (mqtt.Client, error) {

	options := mqtt.NewClientOptions()
	options.AddBroker(brokerURL)

	options.SetClientID(clientID)

	options.SetDefaultPublishHandler(defaultMessagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler

	client, error := CreateClient(options)

	return client, error
}

//reconnect after connection lost
func reconnectHandler() {
	//adding

}
