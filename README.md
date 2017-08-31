# Rafty

gRPC powered raft concensus framework. 

## Build
`$ go build` (obviously).

## Usage

- Start a leader node: `$ ./rafty start --host=127.0.0.1:8000`
- Start a follower: `$ ./rafty join --host=127.0.0.1:8001 --leader=127.0.0.1:8000`

This creates a leader on port 8000, and creates a node on port 8001, connected to the leader. A list of nodes is shares across each nodes. So that in the event of a node failure, a new node can become elected as the leader and have a record of the existing nodes.

