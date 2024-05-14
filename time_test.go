package main

import (
	"fmt"
	"testing"
	"time"
)

func AddDurationToTime(inputTime time.Time, durationString string) (time.Time, error) {
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return time.Time{}, err
	}
	newTime := inputTime.Add(duration)
	return newTime, nil
}
func Aheadif(t time.Time) bool {
	now := time.Now()
	return t.After(now)
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now())
	n := time.Now()
	future, err := AddDurationToTime(n, "2h30m")
	if err != nil {
		panic(err)
	}
	fmt.Println(future)
	fmt.Println(Aheadif(future))

}
