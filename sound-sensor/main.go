package main

import (
	"machine"
	"strconv"
	"time"
)

const (
	threshold           = time.Millisecond * 500
	soundEventThreshold = 20
)

func main() {
	count := 0
	lastSoundEvent := time.UnixMilli(int64(^uint(0) >> 1)) // max

	eventCount := 0
	inputPin := machine.GP28
	inputPin.Configure(machine.PinConfig{Mode: machine.PinInput})

	for {
		if inputPin.Get() {
			if time.Now().Sub(lastSoundEvent) <= threshold {
				eventCount++
			}
			if eventCount >= soundEventThreshold {
				println("sound detected " + strconv.Itoa(count))
				count++
				eventCount = 0
			}
			lastSoundEvent = time.Now()
		}
	}
}

// Ideas:
// if we take 50 samples and 50% of them are HIGH, then we have a sound. loop has no sleeps at all. We would need a way to reset the counter when a threshold has expired
// How many events per second do we get?
// time.Tick every 5 seconds print a map<unixSeconds, events in that second>
