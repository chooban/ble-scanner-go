package main

import (
	"context"
	"log"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

type BtScanner struct {
	atcScanner ATCScanner
	ctx        context.Context
}

func (s BtScanner) RunAtcScanner() chan TemperatureData {
	return s.atcScanner.Run(s.ctx)
}

func NewBtScanner(ctx context.Context) BtScanner {
	device, err := linux.NewDevice()
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(device)
	c := ble.WithSigHandler(context.WithCancel(ctx))
	s := BtScanner{ctx: c, atcScanner: ATCScanner{}}

	return s
}
