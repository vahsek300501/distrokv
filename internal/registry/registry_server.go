package registry

import (
	"context"
	"log/slog"
	"net"

	pb "github.com/Vahsek/distrokv/pkg/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (registryServer *server) RegisterNode(ctx context.Context, request *pb.RegisterNodeRequest) (*pb.RegisterNodeResponse, error) {
	logger := registryServer.logger
	logger.Info("Request to regsister new node")
	err := registryServer.nodeRegistry.RegisterNewNode(request)
	if err != nil {
		logger.Error("Failed to register the node")
		return &pb.RegisterNodeResponse{
			Status:  "500",
			Message: "Node Failed to register",
		}, status.Error(codes.Internal, err.Error())
	}
	logger.Info("Successfully registered the node")
	return &pb.RegisterNodeResponse{
		Status:  "200",
		Message: "Node registered successfully",
	}, nil
}

func (registryServer *server) GetPrimaryNode(ctx context.Context, request *pb.PrimaryNodeRequest) (*pb.PrimaryNodeResponse, error) {
	return nil, nil
}

func (registryServer *server) NodeHeartBeat(ctx context.Context, request *pb.HeartBeatRequest) (*pb.HeartBeatResponse, error) {
	logger := registryServer.logger
	err := registryServer.nodeRegistry.RegisterNodeHeartBeat(request)
	if err != nil {
		logger.Error("Failed to register node heatbeat")
		return &pb.HeartBeatResponse{
			Status:  "500",
			Message: "Failed to register heartbeat",
		}, status.Error(codes.Internal, err.Error())
	}

	logger.Info("successfully registered the node heartbeat")
	return &pb.HeartBeatResponse{
		Status:  "200",
		Message: "Heartbeat registed successfully",
	}, nil
}

func (registryServer *server) GetNodeList(ctx context.Context, request *pb.NodeListRequest) (*pb.NodeListResponse, error) {
	nodeList := registryServer.nodeRegistry.GetNodeList()
	return &pb.NodeListResponse{
		NodeList: nodeList,
	}, nil
}

func StartRegistryServer(portNumber string, logger slog.Logger) {
	logger.Info("Creating TCP Socket on port" + portNumber)
	lis, err := net.Listen("tcp", portNumber)
	if err != nil {
		logger.Error("Error in Creating TCP socket")
	}
	regServer := grpc.NewServer()
	logger.Info("Initializing GRPC service for registry")
	pb.RegisterRegistryServiceServer(regServer, InitializeNewServer(logger))
	if err := regServer.Serve(lis); err != nil {
		logger.Info("Failed to initialize GRPC server for registry")
	} else {
		logger.Info("Successfully initialized GRPC server")
	}
}
