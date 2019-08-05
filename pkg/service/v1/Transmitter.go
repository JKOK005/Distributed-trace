package v1

import (
	pb "Distributed-trace/pkg/api/proto"
	"github.com/golang/protobuf/jsonpb"
	"log"
)

type TransmitterNode struct {}

var (
	kafkaClient *Kclient
	jsonMarshaler = jsonpb.Marshaler{}
)

func (t TransmitterNode) protoToString(msg *pb.TraceReport) (string, error) {
	if data, err := jsonMarshaler.MarshalToString(msg); err != nil {
		return "", err
	}else {return data, nil}
}

func (t TransmitterNode) dispatch(payload *pb.TraceReport) error {
	/*
		Attempts to send received msgs to Kafka
	*/
	if str, err := t.protoToString(payload); err != nil {
		return err
	}else {return kafkaClient.dispatch([]byte(str))}
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