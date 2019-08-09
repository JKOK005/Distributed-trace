### Distributed trace

#### Dependencies
Install [go Kafka-client](https://github.com/confluentinc/confluent-kafka-go) from Confluent 
- Follow instructions listed in the link

Install [go Zookeeper](https://godoc.org/github.com/samuel/go-zookeeper/zk)
- ```go get github.com/samuel/go-zookeeper/zk```

Install google protobuf for go and OS
- protobuf [go client package](https://github.com/golang/protobuf)
- protobuf for Mac OS (google this :/)

Compile protobuf schema via: 
```go
protoc -I api/proto/v1/ --go_out=plugins=grpc:pkg/api/proto/ api/proto/v1/messages.proto
```

### Execution
To start the execution, run 
```dockerfile
docker-compose -f docker-compose.yml up
```

Docker compose will spin up 2 heart beat nodes, 1 ZK node and 1 Kafka standalone node containing 1 broker.

To mine the pinged data, we can spawn a consumer that listens to the Kafka container under the topic _distributedTrace_ 