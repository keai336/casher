package main

import (
	"time"
)

func main() {
	casher := NewOneCasher("test.yaml")
	defer casher.OffDuty()
	for {
		casher.Update()
		time.Sleep(time.Duration(casher.Interval) * time.Second)
	}
}
