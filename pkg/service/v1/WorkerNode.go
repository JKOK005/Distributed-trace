package v1

import (
	"log"
	zk "github.com/samuel/go-zookeeper/zk"
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


func (s SdClient) registerNode(wn WorkerNode) error {
	/* Creates node as ephemeral to ZK cluster under root path */
	log.Println("Registering node address at ", wn.My_address)

	full_path := root_path_zk + "/" + wn.My_address
	_, err := s.conn.CreateProtectedEphemeralSequential(full_path, []byte("Worker"), zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	return nil
}

//func (s SdClient) getNodes() (error) {
//	/* Gets all nodes within path */
//
//}

func (wn WorkerNode) NewClient() (*SdClient, error) {
	// Registers node with ZK cluster
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

func (wn WorkerNode) Start() error {
	client, err := wn.NewClient()
	if err != nil {
		panic(err)
	}

	if err := client.registerNode(wn); err != nil {
		panic(err)
	}

	return nil
}
