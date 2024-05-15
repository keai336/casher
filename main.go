package main

import "time"

func main() {
	casher := NewOneCasher("test.yaml")
	defer casher.OffDuty()
	casher.Update()
	for {
		casher.Update()
		time.Sleep(180 * time.Second)
	}
}
