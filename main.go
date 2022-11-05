package main

import (
	"context"
	"fmt"
)

func main() {
	fmt.Println("Starting BLE scanner")

	config := ReadConfiguration()
	ctx := context.Background()

	scanner := NewBtScanner(ctx)
	ch := scanner.RunAtcScanner()

	mqttClient := NewMQTTClient(&config.Mqtt)
	mqttClient.Run()

	for msg := range ch {
		mqttClient.Publish(&msg)
	}
	fmt.Println("Exiting...")
}
