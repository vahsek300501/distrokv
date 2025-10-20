package clientcommon

import (
	"context"
	"log/slog"
	"sync"

	"google.golang.org/grpc"
)

type ClientFactory struct {
	ConnectionMap map[string]*grpc.ClientConn
	mu            sync.RWMutex
	logger        slog.Logger
}

func InitializeClientFactory(logger slog.Logger) *ClientFactory {
	return &ClientFactory{
		ConnectionMap: make(map[string]*grpc.ClientConn),
		logger:        logger,
	}
}

func (factory *ClientFactory) GetOrCreateConnection(ctx context.Context, builder GrpcBuilder) (*grpc.ClientConn, error) {
	factory.logger.Info("Getting or creating client connection")
	factory.mu.RLock()
	conn, exist := factory.ConnectionMap[builder.target]
	factory.mu.RUnlock()

	if exist && conn != nil {
		factory.logger.Info("Found existing connection")
		return conn, nil
	}

	factory.mu.Lock()
	defer factory.mu.Unlock()
	factory.logger.Info("No exisiting client connection found. Creating a new client connection")

	newConnection, err := builder.Build()
	if err != nil {
		factory.logger.Error("Error in creating the client connection")
		return nil, err
	}
	factory.ConnectionMap[builder.target] = newConnection
	factory.logger.Info("Successfully created the client connection")
	return newConnection, nil
}

func GetClient[T any](ctx context.Context, factory *ClientFactory, builder GrpcBuilder, clientConstructor func(*grpc.ClientConn) T) (T, error) {
	conn, err := factory.GetOrCreateConnection(ctx, builder)
	if err != nil {
		var zero T
		return zero, err
	}
	return clientConstructor(conn), nil
}
