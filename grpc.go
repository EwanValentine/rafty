package main

import (
	pb "github.com/ewanvalentine/rafty/proto"
	"golang.org/x/net/context"
)

type server struct {
	rafty *Rafty
}

func (s *server) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	node, err := s.rafty.AddNode(Node{Host: req.Host})
	return &pb.JoinResponse{Id: node.ID}, err
}

func (s *server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	return &pb.ListResponse{}, nil
}

func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.rafty.mutex.Lock()
	s.rafty.Timeout = TimeoutThreshold
	// Sync data/logs here
	s.rafty.mutex.Unlock()

	return &pb.HeartbeatResponse{Success: true}, nil
}

func (s *server) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	return &pb.RequestVoteResponse{Vote: true}, nil
}
