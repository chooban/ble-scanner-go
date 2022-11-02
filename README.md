# BLE Temperature Scanner

A command for reading Bluetooth Low Energy advertisements from temperature
sensors that have been (https://github.com/pvvx/ATC_MiThermometer)[update with
custom firmware]. Data is then sent to MQTT where it can be distributed,
stored, graphed, whatever.

## Motivation

Once it gets towards winter the temptation is to put on the heating, but we I
need to? Turns out I'm really bad at estimating the temperature, especially if
it's "room temperature". I bought one small temperature sensor which told me
the rooms were warm, so I didn't trust it and bought another. But this one was
flashable! 

## Build

`make build` will yield a binary suitable for running on an RPI. I found it
next to impossible to get a mortal user to be able to see bluetooth data on the
RPI so I've ended up running it as root.

`make debbuild` will yield a dpkg which can be installed which will give you a
systemd service to control.

## Bugs

Probably a fair few. This is my second attempt at using Go to do something, so
it is also probably non-idiomatic and over-engineered.
