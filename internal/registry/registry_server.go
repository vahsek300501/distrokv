package registry

import (
	"context"
	"errors"
	"net"

	registrycontroller "github.com/Vahsek/distrokv/internal/registry/controllers"
	pb "github.com/Vahsek/distrokv/pkg/registry"
	"google.golang.org/grpc"
)

func (registryServer *server) RegisterNode(ctx context.Context, request *pb.RegisterNodeRequest) (*pb.RegisterNodeResponse, error) {
	status := registrycontroller.RegisterNewNode(request)
	if status {
		return &pb.RegisterNodeResponse{
			Status:  "200",
			Message: "Node registered successfully",
		}, nil
	}

	return &pb.RegisterNodeResponse{
		Status:  "500",
		Message: "Node Failed to register",
	}, errors.New("Failed to register node")
}

func (registryServer *server) GetPrimaryNode(ctx context.Context, request *pb.PrimaryNodeRequest) (*pb.PrimaryNodeResponse, error) {
	return nil, nil
}

func (registryServer *server) NodeHeartBeat(ctx context.Context, request *pb.HeartBeatRequest) (*pb.HeartBeatResponse, error) {
	status := registrycontroller.RegisterNodeHeartBeat(request)
	if status {
		return &pb.HeartBeatResponse{
			Status:  "200",
			Message: "Heartbeat registed successfully",
		}, nil
	}
	return &pb.HeartBeatResponse{
		Status:  "500",
		Message: "Failed to register heartbeat",
	}, errors.New("Failed to register heartbeat")
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
