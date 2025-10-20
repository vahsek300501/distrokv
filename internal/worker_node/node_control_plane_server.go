package workernode

import (
	"context"
	"log/slog"
	"net"

	pb "github.com/Vahsek/distrokv/pkg/node/controlplane"
	"google.golang.org/grpc"
)

func (controlPlaneServer *NodeControlPlaneServer) ReplicateSetRequest(ctx context.Context, request *pb.SetReplicationRequest) (*pb.SetReplicationResponse, error) {
	return nil, nil
}

func (controlPlaneServer *NodeControlPlaneServer) ReplicateDeleteRequest(ctx context.Context, request *pb.DeleteReplicationRequest) (*pb.DeleteReplicationResponse, error) {
	return nil, nil
}

func StartNodeControlPlaneServer(controlPlanePortNumber string, logger slog.Logger) {
	logger.Info("Creating TCP Socket on port" + controlPlanePortNumber)
	lis, err := net.Listen("tcp", controlPlanePortNumber)
	if err != nil {
		logger.Error("Error in Creating TCP socket")
	}
	nodeCPServer := grpc.NewServer()
	logger.Info("Initializing GRPC service for node control plane")
	pb.RegisterNodeControlPlaneServiceServer(nodeCPServer, InitializeControlPlaneServer(logger))
	if err := nodeCPServer.Serve(lis); err != nil {
		logger.Info("Failed to initialize GRPC server for node control plane")
	} else {
		logger.Info("Successfully initialized GRPC server")
	}
}
