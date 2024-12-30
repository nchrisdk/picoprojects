package main

import (
	"encoding/json"
	"fmt"
	"log"
	"machine"
	"time"
)

type temp struct {
	TempC float64 `json:"tempC"`
	TempF float64 `json:"tempF"`
}

func main() {

	for {

		t := getTemperature()

		payload, err := json.Marshal(t)
		if err != nil {
			log.Println("error marshalling temperature. Got %v", err)
		}
		fmt.Printf("%s\n", string(payload))
		time.Sleep(time.Second * 5)

	}
}

func getTemperature() *temp {
	curTemp := machine.ReadTemperature()

	return &temp{
		TempC: float64(curTemp) / 1000,
		TempF: ((float64(curTemp) / 1000) * 9 / 5) + 32,
	}
}
