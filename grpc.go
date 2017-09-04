package main

import (
	pb "github.com/ewanvalentine/rafty/proto"
	"golang.org/x/net/context"
)

type server struct {
	rafty *Rafty
}

func (s *server) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	node, err := s.rafty.AddNode(Node{ID: req.Id, Host: req.Host})
	return &pb.JoinResponse{Id: node.ID}, err
}

func (s *server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	return &pb.ListResponse{}, nil
}

func convertProtoNodesToNodes(pbNodes []*pb.Node) []Node {
	var nodes []Node
	for _, node := range pbNodes {
		nodes = append(nodes, Node{
			ID:     node.Id,
			Host:   node.Host,
			Status: node.Status,
		})
	}
	return nodes
}

func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.rafty.mutex.Lock()
	s.rafty.Timeout = TimeoutThreshold
	s.rafty.Leader = req.Leader

	// Here we need to tell the follower about
	// all of the other nodes
	s.rafty.Nodes = convertProtoNodesToNodes(req.Nodes)

	// Sync data/logs here
	s.rafty.mutex.Unlock()

	return &pb.HeartbeatResponse{Success: true}, nil
}

func (s *server) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	return &pb.RequestVoteResponse{Vote: true}, nil
}

func (s *server) AnnounceLeader(ctx context.Context, req *pb.AnnounceLeaderRequest) (*pb.AnnounceLeaderResponse, error) {
	s.rafty.mutex.Lock()
	s.rafty.Leader = req.Id
	s.rafty.Status = FollowerStatus
	s.rafty.Timeout = TimeoutThreshold

	// Establish new connection with the leader
	s.rafty.mutex.Unlock()
	s.rafty.RegisterFollower()

	return &pb.AnnounceLeaderResponse{Consent: true}, nil
}
