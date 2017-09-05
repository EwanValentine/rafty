package main

import (
	"errors"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"golang.org/x/sync/syncmap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/ewanvalentine/rafty/proto"
)

const (
	CandidateStatus = "Candidate"
	FollowerStatus  = "Follower"
	LeaderStatus    = "Leader"

	TimeoutThreshold = 2

	TestMode = "Test"
)

type Attributes struct {
	sm syncmap.Map
}

type Data struct {
	sm syncmap.Map
}

// Node
type Node struct {
	ID   string
	Host string
	Attributes
	Status  string
	Timeout int
	Votes   int
	client  pb.RaftyClient
	conn    *grpc.ClientConn
	mode    string
}

type Rafty struct {
	mutex sync.Mutex
	Node
	Leader       string
	LeaderClient pb.RaftyClient
	LeaderConn   *grpc.ClientConn

	// List of all other Nodes in the network.
	Nodes     []Node
	quit      chan bool
	connected bool

	Data
}

type RaftServer interface {
	Start(host string)
	AddNode(node Node) (Node, error)
	RemoveNode(id string)
	Vote() error
	Timer(election chan<- Node)
}

// Leader - create a node in leader mode
func Leader() *Rafty {
	return &Rafty{
		Nodes: make([]Node, 0),
		Data:  Data{syncmap.Map{}},
		Node: Node{
			ID:      uuid.NewV4().String(),
			Timeout: TimeoutThreshold,
			Votes:   0,
			Status:  LeaderStatus,
		},
	}
}

// Follower - create a node in follower mode
func Follower() *Rafty {
	return &Rafty{
		Nodes: make([]Node, 0),
		Data:  Data{syncmap.Map{}},
		Node: Node{
			ID:      uuid.NewV4().String(),
			Timeout: TimeoutThreshold,
			Votes:   0,
			Status:  FollowerStatus,
		},
	}
}

// Start - start the node, including gRPC server
// also spins up mode specific processes, such as
// the heartbeat process in leader mode
func (rafty *Rafty) Start(host string) {
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Could not start master node: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRaftyServer(s, &server{rafty})
	reflection.Register(s)

	rafty.quit = make(chan bool)

	// Only a leader should perform a heartbeat
	if rafty.Status == LeaderStatus {
		rafty.registerLeader()
	}

	// Only followers should have a timer
	if rafty.Status == FollowerStatus {
		rafty.registerFollower()
	}

	go func() {
		for {
			select {
			case <-rafty.quit:
				log.Println("Gracefully quitting...")
				rafty.connected = false
				s.GracefulStop()
				os.Exit(0)
			}
		}
	}()

	log.Printf(
		"Starting node as %s on host %s",
		rafty.Node.Status,
		host,
	)

	rafty.connected = true

	// Start gRPC server
	if err := s.Serve(lis); err != nil {
		rafty.connected = false
		log.Fatalf("Failed to start master node server: %v", err)
	}
}

// registerLeader - starts the current node in leadership mode
// I.e new status and triggers the heartbeat process to other nodes
func (rafty *Rafty) registerLeader() {
	log.Printf("Node %s leader, performing heartbeat duties \n", rafty.Node.ID)

	// Start REST api to show data
	go API(rafty)

	rafty.RemoveNode(rafty.Node.ID)

	rafty.reconnectAllNodes()

	for _, node := range rafty.Nodes {
		_, err := node.client.AnnounceLeader(
			context.Background(),
			&pb.AnnounceLeaderRequest{Id: rafty.Node.ID},
		)
		log.Println(
			"Failed to inform node %s of leadership: %v",
			node.ID,
			err,
		)
	}

	id := rafty.Node.ID
	if id == "" {
		id = uuid.NewV4().String()
	}

	rafty.mutex.Lock()
	rafty.Leader = id
	rafty.mutex.Unlock()

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

// registerFollower - starts node in follower mode, this starts
// a listener for heartbeats from the leader more
func (rafty *Rafty) registerFollower() {
	log.Println("Node is a follower, listening for heartbeats")
	election := make(chan Node)
	rafty.Timer(election)
	rafty.electionListener(election)
}

// Timer - this starts the countdown timer for
// follower nodes
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

// electionListener -
func (rafty *Rafty) electionListener(election <-chan Node) {
	go func() {
		for {
			select {
			case node := <-election:
				rafty.Node.Status = CandidateStatus
				err := rafty.startElection(node)
				if err != nil {
					log.Fatal("Failed to start election: %v", err)
				}
			}
		}
	}()
}

// startElection -
func (rafty *Rafty) startElection(node Node) error {

	// Vote for self
	err := rafty.Vote()
	if err != nil {
		return err
	}

	rafty.reconnectAllNodes()

	// This needs to be async
	for _, node := range rafty.Nodes {
		if rafty.mode == TestMode {
			rafty.Vote()
			return nil
		}
		resp, err := node.client.RequestVote(
			context.Background(),
			&pb.RequestVoteRequest{Id: node.ID},
		)
		if err != nil {
			return err
		}
		log.Println("Vote recieved")
		if resp.Vote == true {
			rafty.Vote()
		}
	}
	return nil
}

// Join - sets leader node and connects to it
// this method should be used when starting a node
// as a follower to an existing leader.
func (rafty *Rafty) Join(host, leader string) error {
	log.Printf("Connecting to leader: %s", leader)

	// Form a connection to the new node
	conn, err := grpc.Dial(leader, grpc.WithInsecure())
	if err != nil {
		return err
	}

	rafty.LeaderConn = conn
	rafty.LeaderClient = pb.NewRaftyClient(rafty.LeaderConn)

	resp, err := rafty.LeaderClient.Join(
		context.Background(),
		&pb.JoinRequest{Id: rafty.Node.ID, Host: host},
	)
	if err != nil {
		return err
	}

	log.Printf("Joined leader on: %s - %s \n", leader, resp.Id)
	return nil
}

func (rafty *Rafty) isDuplicate(node Node) bool {
	for _, v := range rafty.Nodes {
		if v.ID == node.ID {
			return true
		}
	}
	return false
}

func (rafty *Rafty) isSelf(node Node) bool {
	if rafty.Node.ID == node.ID {
		return true
	}
	return false
}

// AddNode -
func (rafty *Rafty) AddNode(node Node) (Node, error) {
	rafty.mutex.Lock()

	if rafty.mode != TestMode {
		node, err := rafty.connectToNode(node)
		if err != nil {
			return node, err
		}
	}

	if !rafty.isDuplicate(node) && !rafty.isSelf(node) {
		rafty.Nodes = append(rafty.Nodes, node)
	}
	rafty.mutex.Unlock()

	return node, nil
}

// connectToNode -
func (rafty *Rafty) connectToNode(node Node) (Node, error) {

	log.Printf(
		"Connecting to new node %s on %s",
		node.ID,
		node.Host,
	)

	// Form a connection to the new node
	conn, err := grpc.Dial(node.Host, grpc.WithInsecure())
	if err != nil {
		return node, err
	}

	node.conn = conn
	node.client = pb.NewRaftyClient(node.conn)

	return node, nil
}

// reconnectAllNodes -
func (rafty *Rafty) reconnectAllNodes() {
	if rafty.mode == TestMode {
		return
	}
	var nodes []Node
	for _, node := range rafty.Nodes {
		cNode, err := rafty.connectToNode(node)
		if err != nil {

			// If can't connect, node is probably dead,
			// so don't include in list of nodes
			continue
		}
		nodes = append(nodes, cNode)
	}
	rafty.mutex.Lock()
	rafty.Nodes = nodes
	rafty.mutex.Unlock()
}

// RemoveNode -
func (rafty *Rafty) RemoveNode(id string) {
	rafty.mutex.Lock()
	for k, v := range rafty.Nodes {
		if v.ID == id {
			if rafty.mode != TestMode {
				rafty.conn.Close()
			}
			rafty.Nodes = append(rafty.Nodes[:k], rafty.Nodes[k+1:]...)
		}
	}
	rafty.mutex.Unlock()
}

// Vote -
func (rafty *Rafty) Vote() error {
	if rafty.Status != CandidateStatus {
		return errors.New("you cannot vote for a node which is a non-candidate")
	}
	rafty.mutex.Lock()
	rafty.Votes = rafty.Votes + 1
	rafty.mutex.Unlock()

	// This is potentially troubling, if one node dies
	// and an election is performed, this will suffice
	// as we take votes from all nodes, minus the one that died.
	// However, if several nodes die in quick succession, this
	// number may not add up. We need a better way to figure out
	// the majority.
	if rafty.Votes >= len(rafty.Nodes)-1 {
		log.Printf("We have a new leader: %s \n", rafty.Node.ID)
		rafty.mutex.Lock()
		rafty.Status = LeaderStatus
		rafty.mutex.Unlock()
		rafty.registerLeader()
	}

	return nil
}

func convertNodesToProtoNodes(nodes []Node) []*pb.Node {
	var pbNodes []*pb.Node
	for _, v := range nodes {
		pbNodes = append(pbNodes, &pb.Node{
			Id:   v.ID,
			Host: v.Host,
			// @todo - Attributes missing
			Status: v.Status,
		})
	}
	return pbNodes
}

func parseData(data Data) []*pb.Data {
	var newData []*pb.Data
	newData := make(
	data.sm.Range(func(key, value interface{}) bool {
		newData = append(newData, &pb.Data{
			Key:   key.(string),
			Value: value.(string),
		})
		return true
	})
	return newData
}

// Heartbeat - poll all connected nodes with data
func (rafty *Rafty) Heartbeat() error {
	log.Println("Connected nodes: ", len(rafty.Nodes))

	// For each follower node
	for _, node := range rafty.Nodes {

		// Convert syncmap data into protobuf data format
		data := parseData(rafty.Data)
		log.Println(data)

		// Send heartbeat to node with meta data
		// and node data
		_, err := node.client.Heartbeat(
			context.Background(),
			&pb.HeartbeatRequest{
				Leader: rafty.Node.ID,
				Data:   data,
				Nodes:  convertNodesToProtoNodes(rafty.Nodes),
			},
		)

		if err != nil {
			rafty.RemoveNode(node.ID)
			log.Printf("Dead node: %d - %v\n", node.ID, err)
		}
	}
	return nil
}

// Commit - commits data to be stores across all nodes
func (rafty *Rafty) Commit(key string, value interface{}) {
	rafty.Data.sm.Store(key, value)
}
