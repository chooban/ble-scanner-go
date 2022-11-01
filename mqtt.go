package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client   mqtt.Client
	config   *MQTTClientConfig
	transmit chan *TemperatureData
}

func NewMQTTClient(config *MQTTClientConfig) *MQTTClient {
	if config.Broker == "" {
		panic("No MQTT broker defined")
	}
	opts := mqtt.NewClientOptions().AddBroker(config.Broker).SetClientID(config.ClientId)

	client := mqtt.NewClient(opts)

	return &MQTTClient{
		client:   client,
		config:   config,
		transmit: make(chan *TemperatureData),
	}
}

func (c *MQTTClient) Publish(t *TemperatureData) {
	c.transmit <- t
}

func (c *MQTTClient) Run() {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Waiting for MQTT messages...")
	go func() {
		for data := range c.transmit {
			t := *data
			c.client.Publish(c.config.Channel, 1, false, toLineProtocol(&t))
		}
	}()
}

// toLineProtocol function converts TemperatureData struct to MQTT Line Protocol message
func toLineProtocol(t *TemperatureData) string {
	return fmt.Sprintf("%s,id=%s temp=%f,hum=%d,batt=%d", "temperature", t.DeviceName, t.Temperature, t.Humidity, t.Battery)
}
