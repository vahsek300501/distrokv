package service

import (
	"log/slog"

	"github.com/Vahsek/distrokv/internal/storage"
	"github.com/Vahsek/distrokv/internal/worker_node/clients"
)

type NodeConfig struct {
}

type NodeService struct {
	ClusterClient *clients.ClusterClient
	Storage       *storage.KeyValueStore
	logger        slog.Logger
}
