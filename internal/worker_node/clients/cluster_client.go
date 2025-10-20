package clients

import (
	"log/slog"

	clientcommon "github.com/Vahsek/distrokv/internal/common/client_common"
)

type ClusterClient struct {
	factory *clientcommon.ClientFactory
	logger  slog.Logger
}

func InitializeClusterClient(logger slog.Logger) *ClusterClient {
	return &ClusterClient{
		factory: clientcommon.InitializeClientFactory(),
		logger:  logger,
	}
}
