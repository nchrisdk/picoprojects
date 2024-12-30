package main

import (
	"machine"
	"time"
)

func main() {
	pirPin := machine.GP17
	pirPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	for {
		if pirPin.Get() {
			println("motion")
		} else {
			println("reset")
		}
		time.Sleep(time.Millisecond * 200)
	}
}
