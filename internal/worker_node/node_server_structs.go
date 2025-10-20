package workernode

import (
	"log/slog"

	pbControlPlane "github.com/Vahsek/distrokv/pkg/node/controlplane"
	pbDataPlane "github.com/Vahsek/distrokv/pkg/node/dataplane"
)

type NodeControlPlaneServer struct {
	pbControlPlane.UnimplementedNodeControlPlaneServiceServer
	logger slog.Logger
}

type NodeDataPlaneServer struct {
	pbDataPlane.UnimplementedNodeKeyValueServiceServer
	logger slog.Logger
}

func InitializeControlPlaneServer(logger slog.Logger) *NodeControlPlaneServer {
	return &NodeControlPlaneServer{
		logger: logger,
	}
}

func InitializeDataPlaneServer(logger slog.Logger) *NodeDataPlaneServer {
	return &NodeDataPlaneServer{
		logger: logger,
	}
}
