package main

import (
	pkg "Distributed-trace/pkg/service/v1"
	com "Distributed-trace/pkg/service/v1"
	"fmt"
	"sync"
)

type Node interface {
	Start() 	error
}

const (
	poll_timeout 	int32 = 5000 // ms
	poll_interval 	int32 = 10000 // ms
)

var (
	node_addr    string
	node_port 	 int
)

func init() {
	node_addr = com.GetEnvStr("REGISTER_PUBLIC_DNS","localhost")
	node_port = com.GetEnvInt("REGISTER_PUBLIC_PORT", 1111)
}

func main() {
	// go run main.go <node-type> <public IP addr>
	var wg sync.WaitGroup

	wg.Add(1)

	go pkg.NodeListener {Address:fmt.Sprintf("%s:%d", "localhost", node_port)}.RegisterListener()

	go pkg.HeartbeatNode {	My_address: node_addr,
							My_port: node_port,
							Poll_timeout: poll_timeout,
							Poll_interval: poll_interval,
							}.Start()

	go pkg.TransmitterNode{}.Start()

	wg.Wait()
}
