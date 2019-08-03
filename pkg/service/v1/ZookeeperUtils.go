package v1

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strings"
)

var (
	root_path_zk 		string 		= "distributed_trace"
	heartbeatnode_path 	string 		= "heart_beat_nodes"
	sink_path 			string 		= "sinks"
	servers_zk 			[]string 	= []string{"localhost:2181"}
	conn_timeout 		int 		= 10
)

type SdClient struct {
	zk_servers 	[]string
	zk_root    	string
	conn      	*zk.Conn
}

type GenericNode interface {}

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

func (s SdClient) fullPath(from_path string) string {
	if from_path == "" {return root_path_zk}
	return fmt.Sprintf("/%s/%s", root_path_zk, from_path)
}

func (s SdClient) getNodePaths(from_path string) ([]string, error) {
	log.Println("GetNodePaths called at", from_path)
	childs, _, err := s.conn.Children(from_path)
	if err != nil {return nil, err}
	return childs, nil
}

func (s SdClient) CheckRelativePathExists(path string) (bool, error) {
	/* Checks if path relative to root_path/node_path exists */
	return s.checkPathExists(s.fullPath(path))
}

func (s SdClient) RegisterNode(client_path string, data []byte) error {
	/* Registers node at client_path with data */
	full_path := s.fullPath(client_path)
	log.Println("Registering node address at", full_path)
	return s.constructNodesInPath(full_path, "/", data)
}

func (s SdClient) RegisterEphemeralNode(client_path string, data []byte) error {
	/*
		Registers ephemeral node at client_path with data
		If intermediate paths do not exist, we simply create them as a permanent node with empty data
	*/
	full_path := s.fullPath(client_path)
	log.Println("Registering worker ephemeral node address at", full_path)

	full_path_without_last_slice := strings.Split(full_path, "/")
	full_path_without_last := strings.Join(full_path_without_last_slice[ : len(full_path_without_last_slice)-1], "/")

	if err := s.constructNodesInPath(full_path_without_last, "/", nil); err != nil {return err}
	if err := s.constructEphemeralNode(full_path, data); err != nil {return err}
	return nil
}

func (s SdClient) GetNodeValues (from_path []string) ([][]byte, error) {
	/* Passes in a list of node paths and returns the value of the node */
	values := [][]byte{}
	for _, child_path := range from_path {
		full_path := s.fullPath(child_path)
		data, _, err := s.conn.Get(full_path)
		if err != nil {return nil, err}
		values = append(values, data)
	}
	return values, nil
}

func (s SdClient) GetHeartBeatNodePaths(from_path string) ([]string, error) {
	/* Gets all immediate heart beat node paths */
	full_path := s.fullPath(fmt.Sprintf("%s/%s", heartbeatnode_path, from_path))
	return s.getNodePaths(full_path)
}

func (s SdClient) GetSinkNodePaths(from_path string) ([]string, error) {
	/* Gets all immediate sink nodes paths */
	full_path := s.fullPath(fmt.Sprintf("%s/%s", sink_path, from_path))
	return s.getNodePaths(full_path)
}