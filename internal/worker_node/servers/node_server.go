package servers

import (
	"context"
	"log/slog"
	"net"

	pb "github.com/Vahsek/distrokv/pkg/node/dataplane"
	"google.golang.org/grpc"
)

func (dataplaneServer *NodeDataPlaneServer) GetKey(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	return nil, nil
}

func (dataplaneServer *NodeDataPlaneServer) SetKey(ctx context.Context, request *pb.SetRequest) (*pb.SetResponse, error) {
	return nil, nil
}

func (dataPlaneServe *NodeControlPlaneServer) DeleteKey(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return nil, nil
}

func StartNodeDataPlaneServer(dataPlanePortNumber string, logger slog.Logger) {
	logger.Info("Creating TCP Socket on port" + dataPlanePortNumber)
	lis, err := net.Listen("tcp", dataPlanePortNumber)
	if err != nil {
		logger.Error("Error in Creating TCP socket")
		return
	}
	nodeDPServer := grpc.NewServer()
	logger.Info("Initializing GRPC service for node data plane")
	pb.RegisterNodeKeyValueServiceServer(nodeDPServer, InitializeDataPlaneServer(logger))
	if err := nodeDPServer.Serve(lis); err != nil {
		logger.Info("Failed to initialize GRPC server for data plane")
	} else {
		logger.Info("Successfully initialized GRPC server")
	}
}
