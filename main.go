package main

import (
	"log"

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
