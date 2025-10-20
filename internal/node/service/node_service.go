package service

import (
	"log/slog"

	"github.com/Vahsek/distrokv/internal/node/clients"
	"github.com/Vahsek/distrokv/internal/node/storage"
)

type NodeConfig struct {
}

type NodeService struct {
	ClusterClient *clients.ClusterClient
	Storage       *storage.KeyValueStore
	logger        slog.Logger
}
