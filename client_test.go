package main

import (
	"testing"
)

var (
	timerThreshold = 3
)

func TestTimerStartsInFollowerMode(t *testing.T) {
	instance := Follower()
	if instance.Node.Timeout == timerThreshold {
		t.Fail()
	}
}

func TestCanAddNode(t *testing.T) {
	leader := Leader()
	leader.mode = TestMode
	_, err := leader.AddNode(Node{ID: "abc123"})

	if err != nil {
		t.Fail()
	}

	if leader.Nodes[0].ID != "abc123" {
		t.Fail()
	}
}

func TestCanRemoveNode(t *testing.T) {
	leader := Leader()
	leader.mode = TestMode
	leader.AddNode(Node{ID: "abc123"})
	leader.RemoveNode("abc123")

	if len(leader.Nodes) > 0 {
		t.Fail()
	}
}

func TestCantAddSelf(t *testing.T) {
	leader := Leader()
	leader.mode = TestMode
	_, err := leader.AddNode(leader.Node)
	if err != nil {
		t.Fail()
	}

	if len(leader.Nodes) > 0 {
		t.Fail()
	}
}

func TestCantAddDuplicate(t *testing.T) {
	leader := Leader()
	leader.mode = TestMode
	node := Node{ID: "abc"}
	leader.AddNode(node)
	leader.AddNode(node)

	if len(leader.Nodes) > 1 {
		t.Fail()
	}
}

func TestCanTriggerElection(t *testing.T) {
	election := make(chan Node)
	follower := Follower()
	follower.Timeout = 0
	follower.Timer(election)
	done := <-election

	if done.ID == "" {
		t.Fail()
	}
}

func TestCanCommitData(t *testing.T) {
	leader := Leader()
	leader.Commit("testing", "123")

	data := leader.Data["testing"]

	if data != "123" {
		t.Fail()
	}
}
