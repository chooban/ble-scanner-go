package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-ble/ble"
	"github.com/pkg/errors"
)

// TemperatureData struct
type TemperatureData struct {
	DeviceName  string
	Temperature float64
	Humidity    int64
	Battery     int64
	Counter     int64
}

func (t *TemperatureData) LineProtocol() string {
	return fmt.Sprintf("%s,id=%s temp=%f,hum=%d,batt=%d", "temperature", t.DeviceName, t.Temperature, t.Humidity, t.Battery)
}

const ATC_UUID = "181a"

func ParseATCMessage(deviceName, msg string) TemperatureData {
	tempToDivide, _ := strconv.ParseInt(msg[12:16], 16, 16)
	temperature := float64(tempToDivide) / 10

	humidity, _ := strconv.ParseInt(msg[16:18], 16, 16)
	battery, _ := strconv.ParseInt(msg[18:20], 16, 8)

	counter, _ := strconv.ParseInt(msg[24:26], 16, 16)

	return TemperatureData{
		deviceName,
		temperature, humidity, battery, counter,
	}
}

func atcAdvFilter(a ble.Advertisement) bool {
	name := a.LocalName()
	hasATCPrefix := len(name) > 0 && strings.HasPrefix(name, "ATC")

	return hasATCPrefix && len(a.ServiceData()) > 0
}

func atcAdvHandler(ch chan TemperatureData) func(a ble.Advertisement) {
	var lastCounter int64 = -1
	return func(a ble.Advertisement) {
		atcUUID, _ := ble.Parse(ATC_UUID)
		serviceData := a.ServiceData()
		for i := 0; i < len(serviceData); i++ {
			sd := serviceData[i]
			if sd.UUID.Equal(atcUUID) {
				encodedData := hex.EncodeToString(sd.Data)
				parsedMessage := ParseATCMessage(a.LocalName(), encodedData)

				if parsedMessage.Counter != lastCounter {
					lastCounter = parsedMessage.Counter
					ch <- parsedMessage
				}
			}
		}
	}
}

type ATCScanner struct{}

// Run method returns a channel for advertisements, and
// starts scanning. Requires a configured ble library
func (a ATCScanner) Run(ctx context.Context) chan TemperatureData {
	ch := make(chan TemperatureData)

	go func() {
		defer close(ch)
		chkErr(ble.Scan(ctx, true, atcAdvHandler(ch), atcAdvFilter))
	}()

	return ch
}

func chkErr(err error) {
	fmt.Println("Checking error")
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
