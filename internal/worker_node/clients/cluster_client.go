package clients

import (
	"log/slog"
	"sync"

	clientcommon "github.com/Vahsek/distrokv/internal/common/client_common"
	nodecommon "github.com/Vahsek/distrokv/internal/common/node_common"
)

type ClusterClient struct {
	factory         *clientcommon.ClientFactory
	registryAddress string
	nodes           map[string]nodecommon.Node
	mu              sync.Mutex
	logger          slog.Logger
}

func InitializeClusterClient(registryAddress string, logger slog.Logger) *ClusterClient {
	return &ClusterClient{
		factory:         clientcommon.InitializeClientFactory(),
		registryAddress: registryAddress,
		nodes:           make(map[string]nodecommon.Node),
		logger:          logger,
	}
}
