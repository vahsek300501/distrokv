package servers

import (
	"log/slog"

	"github.com/Vahsek/distrokv/internal/worker_node/clients"
	"github.com/Vahsek/distrokv/internal/worker_node/data"
	pbControlPlane "github.com/Vahsek/distrokv/pkg/node/controlplane"
	pbDataPlane "github.com/Vahsek/distrokv/pkg/node/dataplane"
)

type NodeControlPlaneServer struct {
	pbControlPlane.UnimplementedNodeControlPlaneServiceServer
	ClusterClient *clients.ClusterClient
	NodeData      *data.NodeData
	logger        slog.Logger
}

type NodeDataPlaneServer struct {
	pbDataPlane.UnimplementedNodeKeyValueServiceServer
	logger slog.Logger
}

func InitializeControlPlaneServer(logger slog.Logger, client *clients.ClusterClient, nodeData *data.NodeData) *NodeControlPlaneServer {
	return &NodeControlPlaneServer{
		ClusterClient: client,
		NodeData:      nodeData,
		logger:        logger,
	}
}

func InitializeDataPlaneServer(logger slog.Logger) *NodeDataPlaneServer {
	return &NodeDataPlaneServer{
		logger: logger,
	}
}
