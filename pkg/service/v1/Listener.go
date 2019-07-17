package v1

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "Distributed-trace/pkg/api/proto"
	"net"
)

type NodeListener struct {
	address string
}

func (nlis NodeListener) PingNode(ctx context.Context, ping_msg *pb.PingMsg) (*pb.PingMsgResp, error) {
	return &pb.PingMsgResp{IsSuccess : true}, nil
}

func (nlis NodeListener) registerListener() {
	/* Spawns grpc listener */
	log.Println("Spawning listener on", nlis.address)

	lis, err := net.Listen("tcp", nlis.address)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWorkerServiceServer(grpcServer, &NodeListener{})
	grpcServer.Serve(lis)
}
