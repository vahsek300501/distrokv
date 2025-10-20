package clients

import (
	"log/slog"
	"sync"

	clientcommon "github.com/Vahsek/distrokv/internal/common/client_common"
)

type PeerNodes struct {
	hostname    string
	ipAddress   string
	controlPort string
}

type ClusterClient struct {
	factory         *clientcommon.ClientFactory
	registryAddress string
	nodes           map[string]PeerNodes
	mu              sync.Mutex
	logger          slog.Logger
}

func InitializeClusterClient(registryAddress string, logger slog.Logger) *ClusterClient {
	return &ClusterClient{
		factory:         clientcommon.InitializeClientFactory(),
		registryAddress: registryAddress,
		nodes:           make(map[string]PeerNodes),
		logger:          logger,
	}
}
