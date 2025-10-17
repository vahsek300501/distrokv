package registry

import (
	"log/slog"

	"github.com/Vahsek/distrokv/internal/registry/controllers"
	pb "github.com/Vahsek/distrokv/pkg/registry"
)

type server struct {
	pb.UnimplementedRegistryServiceServer
	logger       slog.Logger
	nodeRegistry *controllers.NodeRegistry
}

func InitializeNewServer(logger slog.Logger) *server {
	return &server{
		logger:       logger,
		nodeRegistry: controllers.InitializeNodeRegistry(logger),
	}
}
