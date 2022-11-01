package main

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

func main() {
	fmt.Println("Starting BLE scanner")

	config := ReadConfiguration()
	ctx := context.Background()

	scanner := NewBtScanner(ctx, config.duration)
	ch := scanner.RunAtcScanner()

	mqttClient := NewMQTTClient(&config.Mqtt)
	mqttClient.Run()

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	for msg := range ch {
		s.Stop()
		mqttClient.Publish(&msg)
		s.Start()
	}
	s.Stop()
	fmt.Println("Exiting...")
}
