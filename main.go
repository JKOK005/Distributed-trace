package main

import (
	"flag"
	"fmt"
	"log"
	wk "Distributed-trace/pkg/service/v1"
)

type Node interface {
	Register() 	error
	Start() 	error
}

const (
	poll_timeout int32 = 5000 // ms
)

var (
	node_type    *string
	node_addr    *string
)

func validateAllowableNodeTypes(n string) bool {
	switch n {
	case "worker", "transmitter":
		return true
	default:
		return false
	}
}

func assignNodeType(node_type string, poll_timeout int32) Node {
	switch node_type {
	case "worker":
		return wk.WorkerNode{My_address: *node_addr, Poll_timeout: poll_timeout}
	//case "transmitter":
	//	return wk.TransmitterNode{my_address: node_addr}
	default:
		return nil
	}
}

func init() {
	node_type = flag.String("type", "worker", "Defines the node type")
	node_addr = flag.String("address", "localhost:8080", "Public node address")
	flag.Parse()

	if !validateAllowableNodeTypes(*node_type) {
		panic(fmt.Sprintf("Node type not allowed: %s", *node_type))
	}

	log.SetPrefix(fmt.Sprintf("%s (%s): ", *node_addr, *node_type))
	log.Println("Starting node as type", *node_type)
}

func main() {
	// go run main.go <node-type> <public IP addr>
	node := assignNodeType(*node_type, poll_timeout)
	node.Start()
}
