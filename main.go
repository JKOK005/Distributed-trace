package main

import (
	pb "distributed_tracing/pkg/api/proto"
	"flag"
	"fmt"
	"log"
)

type targets []string

type Node interface {
	start()
}

func (i *targets) String() string {
	return "Node targets"
}

func (i *targets) Set(value string) error {
	fmt.Println(value)
	*i = append(*i, value)
	return nil
}

var (
	node_type    *string
	node_addr    *string
	node_targets targets
)

func validateAllowableNodeTypes(n string) bool {
	switch n {
	case "seed", "worker", "transmitter":
		return true
	default:
		return false
	}
}

func assignNodeType(node_type string) Node {
	switch node_type {
	case "seed":
		return SeedNode{my_address: node_addr, target_addresses: node_targets}
	case "worker":
		return WorkerNode{my_address: node_addr, target_addresses: node_targets}
	case "transmitter":
		return TransmitterNode{my_address: node_addr, target_addresses: node_targets}
	}
}

func init() {
	node_type = flag.String("type", "seed", "Defines the node type")
	node_addr = flag.String("address", "localhost:8080", "Public node address")
	flag.Var(&node_targets, "target", "Defines the nodes to propagate pings to, separated by ',' ")
	flag.Parse()

	if !validateAllowableNodeTypes(*node_type) {
		panic(fmt.Sprintf("Node type not allowed: %s", *node_type))
	}

	log.SetPrefix(fmt.Sprintf("%s (%s): ", *node_addr, *node_type))
	log.Println("Starting node as type", *node_type)
}

func main() {
	// go run main.go <node-type> <public IP addr>
	node := assignNodeType(node_type)

}
