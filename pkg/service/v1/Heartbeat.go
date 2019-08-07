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
}

var (
	reportChannel 		= make(chan *pb.TraceReport)
	heartbeatnode_path 	= GetEnvStr("HEART_BEAT_NODE_PATH","heart_beat_nodes")
)

func (wn HeartbeatNode) getFullPath(from_path string) (string) {
	if from_path == "" {return fmt.Sprintf("/%s/%s", root_path_zk, heartbeatnode_path)}
	return fmt.Sprintf("/%s/%s/%s", root_path_zk, heartbeatnode_path, from_path)
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
				reportChannel <- &pb.TraceReport{FromHostAddr: fmt.Sprintf("%s:%d", wn.My_address, wn.My_port),
												ToHostAddr: fmt.Sprintf("%s:%d", node.My_address, node.My_port),
												ResponseTiming: uint32(0),
												IsTransmissionSuccess:false}
				log.Println(fmt.Sprintf("PingNode attempt to %s/%d failed. Reporting status to sink.", node.My_address, node.My_port))
			} else {
				end := time.Now()
				reportChannel <- &pb.TraceReport{FromHostAddr: fmt.Sprintf("%s:%d", wn.My_address, wn.My_port),
												ToHostAddr: fmt.Sprintf("%s:%d", node.My_address, node.My_port),
												ResponseTiming: uint32(end.Nanosecond() - start.Nanosecond()),
												IsTransmissionSuccess:true}
				log.Println(fmt.Sprintf("PingNode attempt to %s/%d succeeded. Reporting status to sink.", node.My_address, node.My_port))
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
	client, err := newZkClient()
	if err != nil {log.Fatal(err)}

	data, _ := json.Marshal(wn)
	if err := client.RegisterEphemeralNode(wn.getFullPath(fmt.Sprintf("%s:%d", wn.My_address, wn.My_port)), data); err != nil {log.Fatal(err)}

	for{
		select {
		case <- time.NewTicker(time.Duration(wn.Poll_interval) * time.Millisecond).C:
			if node_paths, err := client.GetNodePaths(wn.getFullPath("")); err != nil {
				log.Println(err)
			} else {
				for _, node_path := range node_paths {
					if unmarshalled_node, err := client.GetNodeValue(wn.getFullPath(node_path)); err != nil {
						log.Println(err)
					} else {
						marshalled_node, _ := wn.marshalOne(unmarshalled_node)
						go wn.dispatch(marshalled_node)
					}
				}
			}
		}
	}
}
