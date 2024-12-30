package main

import (
	"log"
	"machine"
	"time"
)

func main() {
	led := machine.GP20 // Raspberry Pi Pico
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		log.Println("blink")
		led.High()
		time.Sleep(time.Second / 2)

		led.Low()
		time.Sleep(time.Second / 2)
	}
}
