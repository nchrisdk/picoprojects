package main

import (
	"machine"
	"time"
	"tinygo.org/x/drivers/servo"
)

var (
	pwm      = machine.PWM0
	servoPin = machine.GP0
)

func main() {
	time.Sleep(5 * time.Second)
	println("ready")
	s, err := servo.New(pwm, servoPin)
	if err != nil {
		println(err.Error())
	}
	s.SetMicroseconds(0)
	for {
		println("setting angle 0")
		s.SetMicroseconds(500) // 0.5 ms
		time.Sleep(time.Second)
		println("setting angle 90")
		s.SetMicroseconds(1500) // 1.5 ms
		time.Sleep(2 * time.Second)
		println("setting angle 180")
		s.SetMicroseconds(2500) // 2.5 ms
		time.Sleep(3 * time.Second)
	}

}
