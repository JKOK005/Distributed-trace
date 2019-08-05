package v1

import (
	pb "Distributed-trace/pkg/api/proto"
	"log"
)

type TransmitterNode struct {}

var (
	kafkaClient *Kclient
)

func (t TransmitterNode) dispatch(payload *pb.TraceReport) error {
	/*
		Attempts to send received msgs to Kafka
	*/
	return kafkaClient.Dispatch(payload)
}

func (t TransmitterNode) Start() {
	if kClient, err := newKafkaProducerClient(); err != nil {
		log.Fatal(err)
	}else {kafkaClient = kClient}

	for{
		payload := <- reportChannel
		go t.dispatch(payload)
	}
}