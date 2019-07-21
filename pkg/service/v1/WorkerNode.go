package v1

import (
	"context"
	"encoding/json"
	pb "Distributed-trace/pkg/api/proto"
	"google.golang.org/grpc"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"sync"
	"time"
)

var (
	root_path_zk 	string 		= "/distributed_trace"
	servers_zk 		[]string 	= []string{"localhost:2181"}
	conn_timeout 	int 		= 10
)

type SdClient struct {
	zk_servers 	[]string
	zk_root    	string
	conn      	*zk.Conn
}

type WorkerNode struct {
	My_address 		string
	Poll_timeout 	int32
	Poll_interval 	int32
}

func (s SdClient) constructZkPath(path string) error {
	log.Println("Creating node at ", path)
	_, err := s.conn.Create(path, []byte(""), 0, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {
		return err
	}
	return nil
}

func (s SdClient) checkPathExists(path string) (bool, error) {
	exists, _, err := s.conn.Exists(path)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s SdClient) registerNode(wn *WorkerNode) error {
	/* Creates node as ephemeral to ZK cluster under root path */
	log.Println("Registering node address at", wn.My_address)

	full_path := root_path_zk + "/" + wn.My_address

	data, err := json.Marshal(wn)
	if err != nil {
		return err
	}

	_, err = s.conn.CreateProtectedEphemeralSequential(full_path, data, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	return nil
}

func (s SdClient) getNodesFromRoot(root_path string) ([]*WorkerNode, error) {
	/* Gets all immediate child nodes that are associated with root_path */
	childs, _, err := s.conn.Children(root_path)

	if err != nil {
		return nil, err
	}
	nodes := []*WorkerNode{}
	for _, each_child := range childs {
		child_path := root_path + "/" + each_child
		data, _, err := s.conn.Get(child_path)
		if err != nil {
			return nil, err
		}
		node := new(WorkerNode)
		err = json.Unmarshal(data, node)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (wn WorkerNode) newClient() (*SdClient, error) {
	/* Registers node with ZK cluster */
	log.Println("Registering to ZK cluster under %s", root_path_zk)

	client := new(SdClient)
	client.zk_servers = servers_zk
	client.zk_root = root_path_zk

	conn, _, err := zk.Connect(servers_zk, time.Duration(conn_timeout) * time.Second)
	if err != nil {
		return nil, err
	}

	client.conn = conn

	if exists, err := client.checkPathExists(root_path_zk); err != nil {
		return nil, err
	} else if exists == false {
		if err := client.constructZkPath(root_path_zk); err != nil {
			return nil, err
		}
	}
	return client, nil
}

func (wn WorkerNode) dispatch(node *WorkerNode) error {
	/* Starts communicating with other nodes via exposed grpc endpoints */
	log.Println("Attempting to communicate with: ", node.My_address)

	if conn, err := grpc.Dial(node.My_address, grpc.WithInsecure()); err != nil {
		log.Println(err)
	} else {
		defer conn.Close()
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(wn.Poll_timeout) * time.Millisecond)
		defer cancel()

		client := pb.NewWorkerServiceClient(conn)
		resp, err := client.PingNode(ctx, &pb.PingMsg{HostAddr: wn.My_address})

		if ctx.Err() == context.DeadlineExceeded {
		// Request timed out. Report as timeout.
		log.Println("Request timed out: ", ctx.Err())
		}else {
			// Request succeeded
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Request successful: ", resp)
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

	if err := client.registerNode(&wn); err != nil {
		panic(err)
	}

	go NodeListener{address:wn.My_address}.registerListener()
	time.Sleep(1000000000) // 1 s

	for{
		select {
		case <- time.NewTicker(time.Duration(wn.Poll_interval) * time.Millisecond).C:
			if nodes, err := client.getNodesFromRoot(root_path_zk); err != nil {
				log.Fatal(err)
			} else {
				go wn.dispatchList(nodes)
			}
		}
	}

	return nil
}
