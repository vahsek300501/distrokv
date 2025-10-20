package clients

import (
	"log/slog"

	clientcommon "github.com/Vahsek/distrokv/internal/common/client_common"
	"github.com/Vahsek/distrokv/internal/worker_node/data"
)

type ClusterClient struct {
	factory  *clientcommon.ClientFactory
	nodeData *data.NodeData
	logger   slog.Logger
}

func InitializeClusterClient(nodeData *data.NodeData, logger slog.Logger) *ClusterClient {
	return &ClusterClient{
		factory:  clientcommon.InitializeClientFactory(),
		nodeData: nodeData,
		logger:   logger,
	}
}
