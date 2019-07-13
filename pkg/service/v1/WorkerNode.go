package v1

import (
	"log"
	zk "github.com/samuel/go-zookeeper/zk"
)

const (
	zk_root string 	= "/distributed_trace"
	zk_url 	string 	= "localhost"
	zk_port int 	= 2181
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
	log.Println("Path %s does not exist. Creating one.", path)
	_, err := s.conn.Create(path, []byte(""), 0, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {
		return err
	}
	return nil
}

func (s SdClient) ensureRootPathExists() error {
	exists, _, err := s.conn.Exists(zk_root)
	if err != nil {
		return err
	} else if exists == true {
		s.constructZkPath(zk_root)
	}
	return nil
}

func (i WorkerNode) Register() error {
	// Registers node with ZK cluster
	log.Println("Registering to ZK cluster under %s", zk_root)


	return nil
}

func (i WorkerNode) Start() error {
	return nil
}
