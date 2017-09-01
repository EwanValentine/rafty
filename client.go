package main

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/ewanvalentine/rafty/proto"
)

const (
	Candidate = "Candidate"
	Follower  = "Follower"
	Leader    = "Leader"

	TimeoutThreshold = 2
)

// Node
type Node struct {
	ID         string
	Host       string
	Attributes map[string]string
	Status     string
	Timeout    int
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
	AddNode(node Node) (Node, error)
	Vote() error
	Heartbeat() error
	Timer(election chan<- Node)
	ElectionListener(election <-chan Node)
	StartElection(node Node)
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

	// Only followers should have a timer
	if rafty.Status == Follower {
		election := make(chan Node)
		log.Println("Testing 123")
		rafty.Timer(election)
		rafty.ElectionListener(election)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start master node server: %v", err)
	}

	log.Println("Started master node")
}

// Timer -
func (rafty *Rafty) Timer(election chan<- Node) {
	timer := time.Tick(1 * time.Second)
	go func() {
		for {
			select {
			case <-timer:
				log.Printf("Timer %d \n", rafty.Node.Timeout)
				if rafty.Node.Timeout == 0 {
					election <- rafty.Node
					return
				}
				rafty.Node.Timeout = rafty.Node.Timeout - 1
			}
		}
	}()
}

// ElectionListener -
func (rafty *Rafty) ElectionListener(election <-chan Node) {
	go func() {
		for {
			select {
			case node := <-election:
				rafty.Node.Status = Candidate
				err := rafty.StartElection(node)
				if err != nil {
					log.Fatal("Failed to start election: %v", err)
				}
			}
		}
	}()
}

// StartElection -
func (rafty *Rafty) StartElection(node Node) error {

	// Vote for self
	err := rafty.Vote()
	if err != nil {
		return err
	}

	// This needs to be async
	for _, node := range rafty.Nodes {
		resp, err := node.client.RequestVote(
			context.Background(),
			&pb.RequestVoteRequest{Id: node.ID},
		)
		if err != nil {
			return err
		}
		log.Println("Vote recieved: %b", resp.Vote)
		rafty.Vote()
	}
	return nil
}

// Join -
func (rafty *Rafty) Join(host, leader string) error {
	_, err := rafty.AddNode(Node{Host: leader})
	if err != nil {
		return err
	}

	log.Printf("Connecting to leader: %s", leader)

	resp, err := rafty.Nodes[0].client.Join(context.Background(), &pb.JoinRequest{Host: host})
	if err != nil {
		return err
	}

	log.Printf("Joined leader on: %s - %d", host, resp.Id)
	return nil
}

// AddNode -
func (rafty *Rafty) AddNode(node Node) (Node, error) {

	// This isn't sufficient, we need to update all other nodes as well
	rafty.mutex.Lock()

	// Form a connection to the new node
	conn, err := grpc.Dial(node.Host, grpc.WithInsecure())
	if err != nil {
		return node, err
	}

	node.ID = uuid.NewV4().String()

	// We need to run conn.Close() when that node dies
	// or is removed somehow
	// defer conn.Close()

	node.client = pb.NewRaftyClient(conn)
	rafty.Nodes = append(rafty.Nodes, node)
	rafty.mutex.Unlock()

	return node, nil
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
		log.Printf("We have a new leader: %s \n", rafty.Node.ID)
		rafty.Status = Leader
	}

	rafty.mutex.Unlock()

	return nil
}

// Heartbeat - poll all connected nodes with data
func (rafty *Rafty) Heartbeat() error {
	log.Println("Connected nodes: ", len(rafty.Nodes))
	for _, node := range rafty.Nodes {
		log.Printf("Data: %d", node.ID)
		_, err := node.client.Heartbeat(context.Background(), &pb.HeartbeatRequest{Data: "test"})
		if err != nil {
			log.Printf("Dead node: %d - %v", node.ID, err)
		}
	}
	return nil
}
