package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

// device = flag.String("device", "default", "implementation of ble")
// var dup = flag.Bool("dup", false, "allow duplicate reported")
var du = flag.Duration("du", 5*time.Second, "scanning duration")

func main() {
	flag.Parse()
	fmt.Println("Starting BLE scanner")

	ctx := context.Background()

	scanner := NewBtScanner(ctx, du)
	ch := scanner.RunAtcScanner()

	mqttClient := NewMQTTClient()
	go mqttClient.Run()

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	for msg := range ch {
		s.Stop()
		fmt.Printf("Received messaged: %+v\n", msg)
		mqttClient.Publish(&msg)
		s.Start()
	}
	s.Stop()
	fmt.Println("Exiting...")
}
