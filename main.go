package main

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/ewanvalentine/rafty/proto"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	start         = kingpin.Command("start", "Start a leader node")
	startHostFlag = start.Flag("host", "Host to start node on").String()

	list         = kingpin.Command("list", "List nodes")
	listHostFlag = list.Flag("host", "Host joinress").String()

	join           = kingpin.Command("join", "Join a node")
	joinHostFlag   = join.Flag("host", "Host address").String()
	joinLeaderFlag = join.Flag("leader", "Leader address").String()
)

const (
	Candidate = "Candidate"
	Follower  = "Follower"
	Leader    = "Leader"
)

// Node
type Node struct {
	ID         int32
	Host       string
	Attributes map[string]string
	Status     string
	Timeout    time.Time
	Votes      int
	client     pb.RaftyClient
}

type Rafty struct {
	mutex sync.Mutex
	Node

	// List of all other Nodes in the network.
	Nodes []Node
}

type RaftServer interface {
	Start(host string)
	AddNode(node Node) error
	Vote() error
	Heartbeat() error
}

// Start -
func (rafty *Rafty) Start(host string) {
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Could not start master node: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRaftyServer(s, &server{rafty})
	reflection.Register(s)

	// Only a leader should perform a heartbeat
	if rafty.Status == Leader {
		log.Println("Node is leader, performing heartbeat duties")
		go func() {
			heartbeat := time.Tick(1 * time.Second)
			for {
				select {
				case <-heartbeat:
					err := rafty.Heartbeat()
					if err != nil {
						log.Printf("Heartbeat error: %v", err)
					}
				}
			}
		}()
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start master node server: %v", err)
	}

	log.Println("Started master node")
}

// Join -
func (rafty *Rafty) Join(host, leader string) error {
	err := rafty.AddNode(Node{Host: leader})
	if err != nil {
		return err
	}

	log.Printf("Connecting to leader: %s", leader)

	log.Println(rafty.Nodes[0].client)

	resp, err := rafty.Nodes[0].client.Join(context.Background(), &pb.JoinRequest{Host: host})
	if err != nil {
		return err
	}

	log.Printf("Joined leader on: %s - %d", host, resp.Id)
	return nil
}

// AddNode -
func (rafty *Rafty) AddNode(node Node) error {

	// This isn't sufficient, we need to update all other nodes as well
	rafty.mutex.Lock()

	log.Println(node)

	// Form a connection to the new node
	conn, err := grpc.Dial(node.Host, grpc.WithInsecure())
	if err != nil {
		return err
	}

	// We need to run conn.Close() when that node dies
	// or is removed somehow
	// defer conn.Close()

	node.client = pb.NewRaftyClient(conn)
	rafty.Nodes = append(rafty.Nodes, node)
	rafty.mutex.Unlock()
	return nil
}

// Vote -
func (rafty *Rafty) Vote() error {
	if rafty.Status != Candidate {
		return errors.New("you cannot vote for a node which is a non-candidate")
	}
	rafty.mutex.Lock()
	rafty.Votes = rafty.Votes + 1

	// This is potentially troubling, if one node dies
	// and an election is performed, this will suffice
	// as we take votes from all nodes, minus the one that died.
	// However, if several nodes die in quick succession, this
	// number may not add up. We need a better way to figure out
	// the majority.
	if rafty.Votes >= len(rafty.Nodes)-1 {
		rafty.Status = Leader
	}

	rafty.mutex.Unlock()

	return nil
}

// Heartbeat -
func (rafty *Rafty) Heartbeat() error {
	log.Println("Node count: ", len(rafty.Nodes))
	for _, node := range rafty.Nodes {
		log.Printf("Hello %d", node.ID)
		log.Println(node.client)
		_, err := node.client.Heartbeat(context.Background(), &pb.HeartbeatRequest{Data: "test"})
		if err != nil {
			log.Printf("Dead node: %d - %v", node.ID, err)
		}
	}
	return nil
}

// rafty create -host 127.0.0.1:8989
// rafty join -host 127.0.0.1:8989
// rafty list
func main() {

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Ewan Valentine")
	kingpin.CommandLine.Help = "Raft concencus (sort of)"

	// Parse cli arguments
	switch kingpin.Parse() {
	case "list":
		log.Println("Listing nodes: ")
	case "join":
		log.Printf("Adding node on address: %s", *joinHostFlag)

		rafty := Rafty{}
		rafty.Status = Follower
		rafty.Votes = 0
		rafty.Join(*joinHostFlag, *joinLeaderFlag)
		rafty.Start(*joinHostFlag)

	case "start":

		// Start server
		rafty := Rafty{}
		rafty.Status = Leader
		rafty.Votes = 0
		rafty.Nodes = make([]Node, 0)

		log.Printf("Starting a node on address: %s", *startHostFlag)
		rafty.Start(*startHostFlag)
	}
}

type server struct {
	rafty *Rafty
}

func (s *server) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	err := s.rafty.AddNode(Node{ID: 1, Host: req.Host})
	return &pb.JoinResponse{Id: 1}, err
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
