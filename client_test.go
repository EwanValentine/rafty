package main

import (
	"testing"
	"time"
)

var (
	timerThreshold = 2
)

func TestTimerStartsInFollowerMode(t *testing.T) {
	leader := Leader()
	go leader.Start(":8000")
	instance := Follower()
	go func() {
		instance.Join(":8001", ":8000")
		instance.Start(":8001")
	}()

	time.Sleep(1 * time.Second)

	if instance.Node.Timeout == timerThreshold {
		t.Fail()
	}

	instance.quit <- true
	leader.quit <- true
}

func TestCanAddNode(t *testing.T) {
	leader := Leader()
	instance := Follower()
	go leader.Start(":8000")
	leader.quit <- true

	go func() {
		instance.Join(":8001", ":8000")
		instance.Start(":8001")
	}()

	if leader.Nodes[0].ID != instance.Node.ID {
		t.Fail()
	}
	instance.quit <- true
	leader.quit <- true
}
