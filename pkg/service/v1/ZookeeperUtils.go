package v1

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strings"
)

var (
	root_path_zk 	string 		= "/distributed_trace"
	node_path 		string 		= "nodes"
	servers_zk 		[]string 	= []string{"localhost:2181"}
	conn_timeout 	int 		= 10
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
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
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
	if err = s.constructNode(pathTrace + pathSlice[len(pathSlice) -1], data); err != nil {return err}
	return nil
}

func (s SdClient) registerNode(client_path string, data []byte) error {
	/* Registers node at client_path with data */
	log.Println("Registering node address at", client_path)
	full_path := fmt.Sprintf("%s/%s/%s", root_path_zk, node_path, client_path)
	return s.constructNodesInPath(full_path, "/", data)
}

func (s SdClient) registerEphemeralNode(client_path string, data []byte) error {
	/*
		Registers ephemeral node at client_path with data
		If intermediate paths do not exist, we simply create them as a permanent node with empty data
	*/
	log.Println("Registering worker ephemeral node address at", client_path)
	full_path := fmt.Sprintf("%s/%s/%s", root_path_zk, node_path, client_path)

	full_path_without_last_slice := strings.Split(full_path, "/")
	full_path_without_last := strings.Join(full_path_without_last_slice[ : len(full_path_without_last_slice)-1], "/")

	if err := s.constructNodesInPath(full_path_without_last, "/", nil); err != nil {return err}
	if err := s.constructEphemeralNode(full_path, data); err != nil {return err}
	return nil
}

func (s SdClient) getNodeValues (from_path []string) ([][]byte, error) {
	/* Passes in a list of node paths and returns the value of the node */
	values := [][]byte{}
	for _, child_path := range from_path {
		full_path := fmt.Sprintf("%s/%s/%s", root_path_zk, node_path, child_path)
		data, _, err := s.conn.Get(full_path)
		if err != nil {return nil, err}
		values = append(values, data)
	}
	return values, nil
}

func (s SdClient) getChildrenPaths(from_path string) ([]string, error) {
	/* Gets all immediate child nodes that are associated with root_path/node_path */
	full_path := fmt.Sprintf("%s/%s/%s", root_path_zk, node_path, from_path)
	childs, _, err := s.conn.Children(full_path)
	if err != nil {return nil, err}
	return childs, nil
}