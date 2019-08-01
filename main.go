package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	pkg "Distributed-trace/pkg/service/v1"
)

type Node interface {
	Start() 	error
}

const (
	poll_timeout 	int32 = 5000 // ms
	poll_interval 	int32 = 10000 // ms
)

var (
	node_type    *string
	node_addr    *string
	node_port 	 *int
)

func validateAllowableNodeTypes(n string) bool {
	switch n {
	case "worker", "transmitter":
		return true
	default:
		return false
	}
}

func init() {
	node_type = flag.String("type", "worker", "Defines the node type")
	node_addr = flag.String("address", "localhost", "Public node address")
	node_port = flag.Int("port", 8000, "Public node port")
	flag.Parse()

	if !validateAllowableNodeTypes(*node_type) {
		panic(fmt.Sprintf("Node type not allowed: %s", *node_type))
	}

	log.SetPrefix(fmt.Sprintf("%s (%s:%d): ", *node_addr, *node_type, *node_port))
	log.Println("Starting node as type", *node_type)
}

func main() {
	// go run main.go <node-type> <public IP addr>
	var wg sync.WaitGroup

	wg.Add(1)
	go pkg.NodeListener {Address:fmt.Sprintf("%s:%d", *node_addr, *node_port)}.RegisterListener()
	go pkg.WorkerNode {My_address: *node_addr, My_port: *node_port, Poll_timeout: poll_timeout, Poll_interval: poll_interval}.Start()
	wg.Wait()
}
