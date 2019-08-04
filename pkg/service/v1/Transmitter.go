package v1

import (
	pb "Distributed-trace/pkg/api/proto"
	"log"
)

type TransmitterNode struct {}

func (t TransmitterNode) dispatch(payload *pb.TraceReport) error {
	/*
		Attempts to send received msgs to Kafka
	*/
	log.Println(payload.ResponseTiming)
	return nil
}

func (t TransmitterNode) Start() {
	for{
	 	payload := <- reportChannel
		go t.dispatch(payload)
	}
}