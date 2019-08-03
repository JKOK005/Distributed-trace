package v1

import (
	pb "Distributed-trace/pkg/api/proto"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

type HeartbeatNode struct {
	My_address 		string
	My_port 		int
	Poll_timeout 	int32
	Poll_interval 	int32
	ReportChannel 	chan *pb.TraceReport
}

func (wn HeartbeatNode) marshalOne(data []byte)(*HeartbeatNode, error) {
	node := new(HeartbeatNode)
	err := json.Unmarshal(data, node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (wn HeartbeatNode) marshalAll(datas [][]byte)([]*HeartbeatNode, error) {
	nodes := []*HeartbeatNode{}
	for _, data := range datas {
		if marshalled, err := wn.marshalOne(data); err != nil {
			return nil, err
		}else{nodes = append(nodes, marshalled)}
	}
	return nodes, nil
}

func (wn HeartbeatNode) dispatch(node *HeartbeatNode) error {
	/* Starts communicating with other nodes via exposed grpc endpoints */
	log.Println("Attempting to communicate with: ", node.My_address, node.My_port)

	if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.My_address, node.My_port), grpc.WithInsecure()); err != nil {
		log.Println(err)
	} else {
		defer conn.Close()
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(wn.Poll_timeout) * time.Millisecond)
		defer cancel()

		client := pb.NewWorkerServiceClient(conn)

		start := time.Now()
		_, err := client.PingNode(ctx, &pb.PingMsg{HostAddr: fmt.Sprintf("%s:%d", node.My_address, node.My_port)})

		if ctx.Err() == context.DeadlineExceeded {
		// Request timed out. Report as timeout.
		log.Println("Request timed out: ", ctx.Err())
		}else {
			// Request succeeded
			if err != nil {
				log.Println(err)
			} else {
				end := time.Now()
				log.Println(fmt.Sprintf("Response received in %d ns", end.Nanosecond() - start.Nanosecond()))
			}
		}
	}
	return nil
}

func (wn HeartbeatNode) dispatchList(nodes []*HeartbeatNode) error {
	for _, each_node := range nodes {go wn.dispatch(each_node)}
	return nil
}

func (wn HeartbeatNode) Start() {
	client, err := newClient()
	if err != nil {log.Fatal(err)}

	data, _ := json.Marshal(wn)
	if err := client.RegisterEphemeralNode(fmt.Sprintf("%s:%d", wn.My_address, wn.My_port), data); err != nil {log.Fatal(err)}

	for{
		select {
		case <- time.NewTicker(time.Duration(wn.Poll_interval) * time.Millisecond).C:
			if node_paths, err := client.GetHeartBeatNodePaths(""); err != nil {
				log.Println(err)
			} else {
				if unmarshalled_nodes, err := client.GetNodeValues(node_paths); err != nil {
					log.Println(err)
				} else {
					marshalled_nodes, _ := wn.marshalAll(unmarshalled_nodes)
					go wn.dispatchList(marshalled_nodes)
				}
			}
		}
	}
}
