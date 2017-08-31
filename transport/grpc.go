package transport

import (
	"context"
	"time"

	"github.com/ewanvalentine/rafty/raft"
)

type server struct {
	rafty *raft.Rafty
}

func (s *server) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	node, err := s.rafty.AddNode(Node{Host: req.Host})
	return &pb.JoinResponse{Id: int32(node.ID)}, err
}

func (s *server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	return &pb.ListResponse{}, nil
}

func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.rafty.mutex.Lock()
	s.rafty.Timeout = time.Now()
	// Sync data/logs here
	s.rafty.mutex.Unlock()

	return &pb.HeartbeatResponse{Success: true}, nil
}

func (s *server) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	return &pb.RequestVoteResponse{Vote: true}, nil
}
