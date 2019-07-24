package v1

import (
	pb "Distributed-trace/pkg/api/proto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type WorkerNode struct {
	My_address 		string
	My_port 		int
	Poll_timeout 	int32
	Poll_interval 	int32
}

func (wn WorkerNode) marshalOne(data []byte)(*WorkerNode, error) {
	node := new(WorkerNode)
	err := json.Unmarshal(data, node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (wn WorkerNode) marshalAll(datas [][]byte)([]*WorkerNode, error) {
	nodes := []*WorkerNode{}
	for _, data := range datas {
		if marshalled, err := wn.marshalOne(data); err != nil {
			return nil, err
		}else{nodes = append(nodes, marshalled)}
	}
	return nodes, nil
}

func (wn WorkerNode) newClient() (*SdClient, error) {
	/* Registers node with ZK cluster */
	log.Println("Registering to ZK cluster")

	client := new(SdClient)

	conn, _, err := zk.Connect(servers_zk, time.Duration(conn_timeout) * time.Second)
	if err != nil {
		return nil, err
	}

	client.conn = conn

	if exists, err := client.checkPathExists(fmt.Sprintf("%s/%s", root_path_zk, node_path)); err != nil {
		return nil, err
	} else if exists == false {
		if err := client.registerNode(fmt.Sprintf("%s/%s", root_path_zk, node_path), []byte{}); err != nil {
			return nil, err
		}
	}
	return client, nil
}

func (wn WorkerNode) dispatch(node *WorkerNode) error {
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

func (wn WorkerNode) dispatchList(nodes []*WorkerNode) error {
	for _, each_node := range nodes {go wn.dispatch(each_node)}
	return nil
}

func (wn WorkerNode) Start(wg *sync.WaitGroup) error {
	client, err := wn.newClient()
	if err != nil {
		panic(err)
	}

	data, _ := json.Marshal(wn)

	if err := client.registerEphemeralNode(fmt.Sprintf("%s/%s/%s:%d",
														root_path_zk, node_path, wn.My_address, wn.My_port), data); err != nil {
		log.Fatal(err)
	}

	go NodeListener {address:fmt.Sprintf("%s:%d", wn.My_address, wn.My_port)}.registerListener()
	time.Sleep(100000000)

	for{
		select {
		case <- time.NewTicker(time.Duration(wn.Poll_interval) * time.Millisecond).C:
			if node_paths, err := client.getChildrenNodes(fmt.Sprintf("%s/%s", root_path_zk, node_path)); err != nil {
				log.Fatal(err)
			} else {
				if unmarshalled_nodes, err := client.getNodeValues(node_paths); err != nil {
					log.Fatal(err)
				} else {
					marshalled_nodes, _ := wn.marshalAll(unmarshalled_nodes)
					go wn.dispatchList(marshalled_nodes)
				}
			}
		}
	}
	return nil
}
