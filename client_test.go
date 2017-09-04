package main

import (
	"log"
	"testing"
	"time"
)

var (
	timerThreshold = 3
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
	wait := make(chan bool)
	leader := Leader()
	go leader.Start(":8000")
	instance := Follower()
	go func() {
		for {
			if leader.connected == true {
				wait <- true
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-wait:
				instance.Join(":8001", ":8000")
				instance.Start(":8001")
			}
		}
	}()

	log.Println(leader.Nodes)

	if leader.Nodes[0].ID != instance.Node.ID {
		t.Fail()
	}
	instance.quit <- true
	leader.quit <- true
}
