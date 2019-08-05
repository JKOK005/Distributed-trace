package v1

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strings"
	"time"
)

var (
	root_path_zk 		string 		= getEnvStr("ROOT_PATH_ZK", "distributed_trace")
	sink_path 			string 		= getEnvStr("SINK_PATH", "sinks")
	servers_zk 			[]string 	= getEnvStrSlice("SERVERS_ZK", []string{"localhost:2181"})
	conn_timeout 		int 		= getEnvInt("SINK_PATH", 10)
)

type SdClient struct {
	zk_servers 	[]string
	zk_root    	string
	conn      	*zk.Conn
}

type GenericNode interface {}

func newZkClient() (*SdClient, error) {
	/* Registers node with ZK cluster */
	log.Println("Creating client to ZK server", servers_zk)
	client := new(SdClient)
	conn, _, err := zk.Connect(servers_zk, time.Duration(conn_timeout) * time.Second)
	if err != nil {return nil, err}

	log.Println("Successfully connected to ZK at", servers_zk)
	client.conn = conn
	return client, nil
}

func (s SdClient) checkPathExists(path string) (bool, error) {
	exists, _, err := s.conn.Exists(path)
	if err != nil {return false, err}
	return exists, nil
}

func (s SdClient) constructNode(path string, data []byte) error {
	/*
		Checks if node exists at path
		Else attempts to create a permanent node with data
	*/
	if exists, err := s.checkPathExists(path); err != nil {
		return err
	} else if exists == false {
		log.Println("Attempting to create node at path", path)
		_, err := s.conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (s SdClient) constructEphemeralNode(path string, data []byte) error {
	/*
		Checks if node exists at path
		Else attempts to an ephemeral node with data
	*/
	if exists, err := s.checkPathExists(path); err != nil {
		return err
	} else if exists == false {
		_, err := s.conn.CreateProtectedEphemeralSequential(path, data, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {return err}
	}
	return nil
}

func (s SdClient) constructNodesInPath(path string, delimiter string, data []byte) error {
	/*
		Creates a ZK path of nested nodes from path and delimiter
		If the node created is not the end node, we will populate its data with an empty []byte
		If the node created is the end node, we will populate its data with given data

		Path 		- /distributed_trace/nodes
		Demiliter 	- '/'
	*/
	var err error
	pathSlice := strings.Split(path, delimiter)
	pathTrace := "/"
	for _, eachPath := range pathSlice[ : len(pathSlice)-1] {
		pathTrace = pathTrace + eachPath
		if err = s.constructNode(pathTrace, nil); err != nil {return err}
	}
	if err = s.constructNode(path, data); err != nil {return err}
	return nil
}

func (s SdClient) GetNodePaths(from_path string) ([]string, error) {
	log.Println("GetNodePaths called at", from_path)
	childs, _, err := s.conn.Children(from_path)
	if err != nil {return nil, err}
	return childs, nil
}

func (s SdClient) GetNodeValue (from_path string) ([]byte, error) {
	/* Passes in a list of node paths and returns the value of the node */
	data, _, err := s.conn.Get(from_path)
	if err != nil {return nil, err}
	return data, nil
}

func (s SdClient) CheckRelativePathExists(path string) (bool, error) {
	/* Checks if path relative to root_path/node_path exists */
	return s.checkPathExists(path)
}

func (s SdClient) RegisterNode(client_path string, data []byte) error {
	/* Registers node at client_path with data */
	log.Println("Registering node address at", client_path)
	return s.constructNodesInPath(client_path, "/", data)
}

func (s SdClient) RegisterEphemeralNode(client_path string, data []byte) error {
	/*
		Registers ephemeral node at client_path with data
		If intermediate paths do not exist, we simply create them as a permanent node with empty data
	*/
	log.Println("Registering worker ephemeral node address at", client_path)

	full_path_without_last_slice := strings.Split(client_path, "/")
	full_path_without_last := strings.Join(full_path_without_last_slice[ : len(full_path_without_last_slice)-1], "/")

	if err := s.constructNodesInPath(full_path_without_last, "/", nil); err != nil {return err}
	if err := s.constructEphemeralNode(client_path, data); err != nil {return err}
	return nil
}