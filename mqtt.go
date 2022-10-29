package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client   mqtt.Client
	transmit chan *TemperatureData
}

func NewMQTTClient() *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker("192.168.1.243:1883").SetClientID("pihole")

	client := mqtt.NewClient(opts)

	return &MQTTClient{
		client:   client,
		transmit: make(chan *TemperatureData),
	}
}

func (c *MQTTClient) Publish(t *TemperatureData) {
	fmt.Println("Publishing message to channel")
	c.transmit <- t
}

func (c *MQTTClient) Run() {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Waiting for MQTT messages...")
	for data := range c.transmit {
		t := *data
		fmt.Printf("To publish: %s\n", t.LineProtocol())
		c.client.Publish("sensor/temperature", 1, false, t.LineProtocol())
	}
}
