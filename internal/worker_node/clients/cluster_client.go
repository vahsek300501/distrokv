package clients

import (
	"context"
	"log/slog"

	clientcommon "github.com/Vahsek/distrokv/internal/common/client_common"
	pb_node_control_plane "github.com/Vahsek/distrokv/pkg/node/controlplane"
	pb_registry "github.com/Vahsek/distrokv/pkg/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

// createRegistryClient creates a gRPC client for registry communication
func (clusterClient *ClusterClient) createRegistryClient(registryServerAddress string) (pb_registry.RegistryServiceClient, error) {
	grpcDialOption := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	registryBuilder := clientcommon.NewGrpcBuilder(registryServerAddress).
		SetGrpcDialOptions(grpcDialOption...)
	registryClientConstructor := func(conn *grpc.ClientConn) pb_registry.RegistryServiceClient {
		return pb_registry.NewRegistryServiceClient(conn)
	}
	client, err := clientcommon.GetClient(
		context.Background(),
		clusterClient.factory,
		*registryBuilder,
		registryClientConstructor)

	if err != nil {
		clusterClient.logger.Error("Error in creating registry client")
		return nil, err
	}
	return client, nil
}

func (clusterClient *ClusterClient) createPeerClientConnection(peerAddress string) (pb_node_control_plane.NodeControlPlaneServiceClient, error) {
	grpcDialOption := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	peerClientBuilder := clientcommon.NewGrpcBuilder(peerAddress).
		SetGrpcDialOptions(grpcDialOption...)
	peerClientConstructor := func(conn *grpc.ClientConn) pb_node_control_plane.NodeControlPlaneServiceClient {
		return pb_node_control_plane.NewNodeControlPlaneServiceClient(conn)
	}
	client, err := clientcommon.GetClient(
		context.Background(),
		clusterClient.factory,
		*peerClientBuilder,
		peerClientConstructor)
	if err != nil {
		clusterClient.logger.Error("Error in creating peer client")
		return nil, err
	}
	return client, nil
}
