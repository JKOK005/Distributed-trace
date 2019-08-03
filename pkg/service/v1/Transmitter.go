package v1

import (
	"math/rand"
	pb "Distributed-trace/pkg/api/proto"
)

type TransmitterNode struct {
	ReportChannel 	chan *pb.TraceReport
}

func (t TransmitterNode) getSinkNodeAddrs() ([]string, error) {
	/*
		Gets a Sink Node Address dynamically from ZK cluster
	*/
	client, err := newClient()
	if err != nil {return nil, err}
	if node_paths, err := client.GetNodePaths(""); err == nil {
		return node_paths, nil
	}else {return nil, err}
}

func (t TransmitterNode) selectAddr(addrs []string) (string) {
	/*
		Selects a sink node address out of all addresses
	*/
	return addrs[rand.Intn(len(addrs))]
}

func (t TransmitterNode) relayToSink(report *pb.TraceReport) error {
	/*
		Attempts to send received msgs to Sink node
		Sink node must be registered in ZK for discovery
	*/
	return nil
}

func (t TransmitterNode) Start() {

}