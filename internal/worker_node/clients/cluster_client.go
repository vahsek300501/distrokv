package clients

import (
	"fmt"
	"log/slog"

	clientcommon "github.com/Vahsek/distrokv/internal/common/client_common"
	"github.com/Vahsek/distrokv/internal/worker_node/data"
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
func (clusterClient *ClusterClient) createRegistryClient(nodeData *data.NodeData) (pb_registry.RegistryServiceClient, error) {
	grpcDialOption := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	registryBuilder := clientcommon.NewGrpcBuilder(nodeData.RegistryServerAddress).
		SetGrpcDialOptions(grpcDialOption...)

	conn, err := registryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create registry connection: %w", err)
	}

	client := pb_registry.NewRegistryServiceClient(conn)
	return client, nil
}

func (clusterClient *ClusterClient) createPeerClientConnection(peerAddress string) (pb_node_control_plane.NodeControlPlaneServiceClient, error) {
	grpcDialOption := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	peerClientBuilder := clientcommon.NewGrpcBuilder(peerAddress).
		SetGrpcDialOptions(grpcDialOption...)

	conn, err := peerClientBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	client := pb_node_control_plane.NewNodeControlPlaneServiceClient(conn)
	return client, nil
}
