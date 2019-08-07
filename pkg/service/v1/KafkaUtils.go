package v1

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"time"
)

var (
	kafka_topic 			= GetEnvStr("KAFKA_TOPIC", "distributedTrace")
	kafka_producer_group 	= GetEnvStr("KAFKA_PRODUCER_GROUP", "distributed_trace_grp")
	kafka_bootstrap_servers = GetEnvStr("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
)

type Kclient struct {
	conn *kafka.Producer
}

func newKafkaProducerClient() (*Kclient, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers" : kafka_bootstrap_servers,
		"group.id" : kafka_producer_group,
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