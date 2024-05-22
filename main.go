package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var path string = "config.yaml"
	args := os.Args
	if len(args) >= 2 {
		path = args[1]
	}
	fmt.Println("path:", path)
	casher := NewOneCasher(path)
	go func() {
		for {
			casher.Update()
			time.Sleep(time.Duration(casher.Interval) * time.Second)
		}

	}()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	casher.OffDuty()
	os.Exit(1)
}
