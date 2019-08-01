package v1

import (
	pb "Distributed-trace/pkg/api/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type NodeListener struct {
	Address string
}

func (nlis NodeListener) PingNode(ctx context.Context, ping_msg *pb.PingMsg) (*pb.PingMsgResp, error) {
	return &pb.PingMsgResp{IsSuccess : true}, nil
}

func (nlis NodeListener) RegisterListener() {
	/* Spawns grpc listener */
	log.Println("Spawning listener on", nlis.Address)

	lis, err := net.Listen("tcp", nlis.Address)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWorkerServiceServer(grpcServer, &NodeListener{})
	grpcServer.Serve(lis)
}
