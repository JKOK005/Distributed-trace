package v1

import (
	"github.com/golang/protobuf/proto"
	pb "Distributed-trace/pkg/api/proto"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"time"
)

var (
	kafka_topic 			= "distributedTrace"
	kafka_consumer_group 	= "distributed_trace_grp"
	kafka_bootstrap_servers = "localhost:9092"
)

type Kclient struct {
	conn *kafka.Producer
}

func newKafkaProducerClient() (*Kclient, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers" : kafka_bootstrap_servers,
		"group.id" : kafka_consumer_group,
	})
	if err != nil {return nil, err}
	return &Kclient{conn: p}, nil
}

func (k Kclient) dispatch(payload []byte) error {
	return k.conn.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &kafka_topic, Partition: kafka.PartitionAny},
					Value:          payload,
					Key:            nil,
					Timestamp:      time.Time{},
					TimestampType:  0,
					Opaque:         nil,
					Headers:        nil,
				}, nil)
	}

func (k Kclient) Dispatch(payload *pb.TraceReport) error {
	if data, err := proto.Marshal(payload); err != nil {
		return err
	} else {return k.dispatch(data)}
}