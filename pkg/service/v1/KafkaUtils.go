package v1

import (
	kafka "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	kafka_consumer_group 	= "distributed_trace"
	kafka_bootstrap_servers = "localhost:9092"
)

type Kclient struct {
	conn *kafka.Consumer
}

func newKafkaClient() (*Kclient, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers" : kafka_bootstrap_servers
	})
}