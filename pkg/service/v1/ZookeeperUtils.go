package v1

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
)

var (
	root_path_zk 	string 		= "/distributed_trace"
	node_path 		string 		= "/nodes"
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
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s SdClient) registerNode(client_path string, data []byte) error {
	/* Creates node to ZK cluster under root path */
	log.Println("Registering node address at", client_path)
	full_path := fmt.Sprintf("%s/%s/%s", root_path_zk, node_path, client_path)
	_, err := s.conn.Create(full_path, data, 0, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {
		return err
	}
	return nil
}

func (s SdClient) registerEphemeralNode(client_path string, data []byte) error {
	/* Creates node as ephemeral to ZK cluster under root path */
	log.Println("Registering worker ephemeral node address at", client_path)
	full_path := fmt.Sprintf("%s/%s/%s", root_path_zk, node_path, client_path)
	_, err := s.conn.CreateProtectedEphemeralSequential(full_path, data, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	return nil
}

func (s SdClient) getNodeValues (node_type GenericNode) func (node_paths []string) ([]*GenericNode, error) {
	/* Passes in a list of node paths and returns the value of the node */

	return func(node_paths []string) ([]*GenericNode, error){
		nodes := []*GenericNode{}
		for _, child_path := range node_paths {
			data, _, err := s.conn.Get(child_path)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(data, node_type)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, &node_type)
		}
		return nodes, nil
	}
}

func (s SdClient) getChildrenNodes(parent_path string) ([]string, error) {
	/* Gets all immediate child nodes that are associated with root_path */
	childs, _, err := s.conn.Children(parent_path)

	if err != nil {
		return nil, err
	}
	return childs, nil
}