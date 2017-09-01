package main

import (
	"testing"
	"time"
)

var (
	timerThreshold = 10
)

func getInstance(status string) Rafty {
	rafty := Rafty{}
	rafty.Status = status
	rafty.Votes = 0
	rafty.Nodes = make([]Node, 0)
	rafty.Node.Timeout = timerThreshold
	return rafty
}

func TestTimerStartsInFollowerMode(t *testing.T) {
	instance := getInstance(Follower)
	go instance.Start(":8000")

	time.Sleep(1 * time.Second)

	if instance.Node.Timeout == timerThreshold {
		t.Fail()
	}
}
