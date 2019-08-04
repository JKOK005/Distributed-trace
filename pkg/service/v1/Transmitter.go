package v1

import (
	pb "Distributed-trace/pkg/api/proto"
)

type TransmitterNode struct {}

func (t TransmitterNode) dispatch(*pb.TraceReport) error {
	/*
		Attempts to send received msgs to Kafka
	*/
	return nil
}

func (t TransmitterNode) Start() {
	for{
		select {
		case <- reportChannel:
			payload := <- reportChannel
			go t.dispatch(payload)
		}
	}
}