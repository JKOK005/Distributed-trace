### Distributed trace

#### Dependencies
Install [go Kafka-client](https://github.com/confluentinc/confluent-kafka-go) from Confluent 
- Follow instructions listed in the link

Install [go Zookeeper](https://godoc.org/github.com/samuel/go-zookeeper/zk)
- ```go get github.com/samuel/go-zookeeper/zk```

Install google protobuf for go and OS
- protobuf [go client package](https://github.com/golang/protobuf)
- protobuf for Mac OS (google this :/)

#### Compile protobuf schema
Execute the command: 
```go
protoc -I api/proto/v1/ --go_out=plugins=grpc:pkg/api/proto/ api/proto/v1/messages.proto
```