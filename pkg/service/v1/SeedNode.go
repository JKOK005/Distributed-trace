package v1

import (
	pb "distributed_tracing/pkg/api/proto"
	"net"
)

const node_type string = "seed"

type node struct {
	address      string
	poll_timeout uint32
}

func (n *node) Init(public_address string, timeout uint32) {
	n.address = public_address
	n.poll_timeout = timeout
}
