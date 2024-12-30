package main

import (
	"fmt"
	"machine"
	"time"
	"tinygo.org/x/drivers/hd44780i2c"
)

type temp struct {
	TempC float64 `json:"tempC"`
	TempF float64 `json:"tempF"`
}

func main() {
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{})

	lcd := hd44780i2c.New(machine.I2C0, 0x27)
	lcd.Configure(hd44780i2c.Config{
		Width:  16,
		Height: 2,
	})
	lcd.SetCursor(0, 0)
	lcd.Print([]byte("Current Temp:"))
	for {
		t := getTemperature()
		lcd.SetCursor(0, 1)
		lcd.Print([]byte(fmt.Sprintf(" %.2f C", t.TempC)))
		time.Sleep(time.Second * 2)
	}
}

func getTemperature() *temp {
	curTemp := machine.ReadTemperature()

	return &temp{
		TempC: float64(curTemp) / 1000,
		TempF: ((float64(curTemp) / 1000) * 9 / 5) + 32,
	}
}

func counter() {
	{
		i2c := machine.I2C0
		i2c.Configure(machine.I2CConfig{})

		lcd := hd44780i2c.New(machine.I2C0, 0x27)
		lcd.Configure(hd44780i2c.Config{
			Width:  16,
			Height: 2,
		})
		lcd.SetCursor(0, 0)
		lcd.Print([]byte("Counting: "))
		var counter uint16 = 0
		for true {
			if counter%100 == 0 {
				lcd.SetCursor(0, 1)
				lcd.Print([]byte(fmt.Sprintf("%v", counter)))
			}
			counter += 100
		}

	}
}

func helloworld() {
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{})

	lcd := hd44780i2c.New(machine.I2C0, 0x27)
	lcd.Configure(hd44780i2c.Config{
		Width:  16,
		Height: 2,
	})
	lcd.Print([]byte(" Hello, world\n LCD 16x02"))
}
