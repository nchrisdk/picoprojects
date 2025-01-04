package main

import (
	"fmt"
	"machine"
	"strconv"
	"time"
	"tinygo.org/x/drivers/dht"
	"tinygo.org/x/drivers/hd44780i2c"
)

const (
	dhtPin = machine.GP28
)

type measurement struct {
	temp        float64
	humidity    float64
	internalRef float64
}

// the I2C bus is used two places and if they don't match the LCD will not show anything!!!
var i2c = machine.I2C1
var dhtSensor = dht.New(dhtPin, dht.DHT11)
var lcd = hd44780i2c.New(i2c, 0x27)

func main() {
	err := i2c.Configure(machine.I2CConfig{
		SDA: machine.GP26,
		SCL: machine.GP27,
	})
	if err != nil {
		println(err.Error())
	}

	lcd.Configure(hd44780i2c.Config{
		Width:  16,
		Height: 2,
	})
	lcd.ClearDisplay()
	for {
		m, err := getMeasurement()
		if err != nil {
			println("dht sensor: " + err.Error())
			continue
		}
		tMessage := strconv.FormatFloat(m.temp, 'f', 1, 64) + "C (" + strconv.FormatFloat(m.internalRef, 'f', 1, 64) + ")"
		hMessage := strconv.FormatFloat(m.humidity, 'f', 1, 64) + "%"
		println(tMessage)
		println(hMessage)
		printLine1(tMessage)
		printLine2(hMessage)

		time.Sleep(2 * time.Second) // Measurements cannot be updated only 2 seconds. More frequent calls will return the same value
	}
}

func printLine1(str string) {
	lcd.SetCursor(0, 0)
	lcd.Print([]byte(str))
}

func printLine2(str string) {
	lcd.SetCursor(0, 1)
	lcd.Print([]byte(str))
}

func getMeasurement() (*measurement, error) {
	t, err := dhtSensor.TemperatureFloat(dht.C)
	if err != nil {
		return nil, fmt.Errorf("temperature: %w", err)
	}
	h, err := dhtSensor.HumidityFloat()
	if err != nil {
		return nil, fmt.Errorf("humidity: %w", err)
	}
	curTemp := machine.ReadTemperature()
	return &measurement{
		temp:        float64(t),
		humidity:    float64(h),
		internalRef: float64(curTemp) / 1000,
	}, nil
}
