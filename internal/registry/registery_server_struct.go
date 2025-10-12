package registry

import (
	"log/slog"

	pb "github.com/Vahsek/distrokv/pkg/registry"
)

type server struct {
	pb.UnimplementedRegistryServiceServer
	logger slog.Logger
}

func InitializeNewServer(logger slog.Logger) *server {
	return &server{
		logger: logger,
	}
}
