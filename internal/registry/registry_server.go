package registry

import (
	"context"
	"net"

	pb "github.com/Vahsek/distrokv/pkg/registry"
	"google.golang.org/grpc"
)

func (registryServer *server) RegisterNode(ctx context.Context, request *pb.RegisterNodeRequest) (*pb.RegisterNodeResponse, error) {
	return nil, nil
}

func (registryServer *server) GetPrimaryNode(ctx context.Context, request *pb.PrimaryNodeRequest) (*pb.PrimaryNodeResponse, error) {
	return nil, nil
}

func StartRegistryServer(portNumber string) {
	lis, err := net.Listen("tcp", portNumber)
	if err != nil {
		// TODO add logging
	}
	regServer := grpc.NewServer()
	pb.RegisterRegistryServiceServer(regServer, &server{})
	if err := regServer.Serve(lis); err != nil {
		// TODO add logging
	}
}
