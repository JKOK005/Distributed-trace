## Distributed trace
Distributed trace aims to track inter-node network performance across a cluster of machines.

This is motivated mostly by the nature of distributed workload across a cluster of machines.
A cluster compromises of different worker nodes with a single master node.
Each worker node can be located under different locality zones or have individual network policies enforced around it.
Moreover, each node can have different instance group with different amount of bandwidth.

The implementation of distributed trace closely follows a work done by [Microsoft](https://conferences.sigcomm.org/sigcomm/2015/pdf/papers/p139.pdf).

At present, we deploy a heartbeat agent to each node in the cluster. The agent is responsible for discovering all nodes within the cluster and reporting the time taken to issue pings to all other nodes within the cluster.

## Installation & deployment
#### Ensure that all dependencies are installed 
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

### Deploying locally via docker compose
To start the execution, run 
```shell script
docker-compose -f docker-compose.yml up
```

Docker compose will spin up 2 heart beat nodes, 1 ZK node and 1 Kafka standalone node containing 1 broker.

To mine the pinged data, we can spawn a consumer that listens to the Kafka container under the topic _distributedTrace

### Deploying locally via minikube
Alternatively, we can deploy the application onto a local minikube cluster:

Start minikube server
```shell script
minikube start
```

Start the Confluent-Kafka chart first
```shell script
helm install --name service-kafka helm-charts/kafka
```

If your `--name` flag is set to another value, ensure that you update the value of KAFKA_BOOTSTRAP_SERVERS in heartbeat chart deployment
```shell script
# Modify the value for KAFKA_BOOTSTRAP_SERVERS in heartbeat chart deployment
value:  "<name>.{{ $.Release.Namespace }}.svc.cluster.local:9092"
```

There after, start distributed trace helm chart
``` shell script
helm install --name distributed-trace helm-charts/dt-chart
```

In deployments to Minikube, you will encounter an issue where the pod is unable to ping its own service. 
This is an existing bug in minikube as of version 1.3.1.

Reference can be found [here](https://github.com/kubernetes/minikube/issues/1568) 

The work around is to do this:
```shell script
minikube ssh
sudo ip link set docker0 promisc on
```