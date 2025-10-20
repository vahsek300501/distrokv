package servers

import (
	"context"
	"log/slog"
	"net"

	"github.com/Vahsek/distrokv/internal/worker_node/clients"
	"github.com/Vahsek/distrokv/internal/worker_node/controllers"
	"github.com/Vahsek/distrokv/internal/worker_node/data"
	pb "github.com/Vahsek/distrokv/pkg/node/controlplane"
	"google.golang.org/grpc"
)

func (controlPlaneServer *NodeControlPlaneServer) ReplicateSetRequest(ctx context.Context, request *pb.SetReplicationRequest) (*pb.SetReplicationResponse, error) {
	return nil, nil
}

func (controlPlaneServer *NodeControlPlaneServer) ReplicateDeleteRequest(ctx context.Context, request *pb.DeleteReplicationRequest) (*pb.DeleteReplicationResponse, error) {
	return nil, nil
}

func (controlPlaneServer *NodeControlPlaneServer) RegisterNewPeerServer(ctx context.Context, request *pb.NewServerAddRequest) (*pb.NewServerAddResponse, error) {
	controlPlaneServer.logger.Info("Registeration request from peer")
	controlPlaneServer.logger.Info(request.String())
	err := controllers.RegisterNewPeerNode(request, controlPlaneServer.NodeData, &controlPlaneServer.logger)
	if err != nil {
		controlPlaneServer.logger.Error("Error in registering new peer node")
		return &pb.NewServerAddResponse{
			Status:  "Failed",
			Message: "Failed to register new server",
		}, err
	}
	return &pb.NewServerAddResponse{
		Status:  "Success",
		Message: "Successfully registerd the new peer server",
	}, nil
}

func StartNodeControlPlaneServer(controlPlanePortNumber string, logger slog.Logger, client *clients.ClusterClient, nodeData *data.NodeData) {
	logger.Info("Creating TCP Socket on port" + controlPlanePortNumber)
	lis, err := net.Listen("tcp", controlPlanePortNumber)
	if err != nil {
		logger.Error("Error in Creating TCP socket")
		return
	}
	nodeCPServer := grpc.NewServer()
	logger.Info("Initializing GRPC service for node control plane")
	pb.RegisterNodeControlPlaneServiceServer(nodeCPServer, InitializeControlPlaneServer(logger, client, nodeData))
	if err := nodeCPServer.Serve(lis); err != nil {
		logger.Info("Failed to initialize GRPC server for node control plane")
	} else {
		logger.Info("Successfully initialized GRPC server")
	}
}
