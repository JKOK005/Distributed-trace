package v1

import (
	"fmt"
	"log"
	"github.com/samuel/go-zookeeper/zk"
	pb "Distributed-trace/pkg/api/proto"
)

type TransmitterNode struct {
	ReportChannel 	chan *pb.TraceReport
}

func (t TransmitterNode) getRandomSinkNodeAddr() (string, error) {
	/*
		Gets a Sink Node Address dynamically from ZK cluster
	*/
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