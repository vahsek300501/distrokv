package clients

import (
	"context"
	"fmt"
	"time"

	clientcommon "github.com/Vahsek/distrokv/internal/common/client_common"
	nodecommon "github.com/Vahsek/distrokv/internal/common/node_common"
	pb_registry "github.com/Vahsek/distrokv/pkg/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// createRegistryClient creates a gRPC client for registry communication
func (clusterClient *ClusterClient) createRegistryClient() (pb_registry.RegistryServiceClient, *grpc.ClientConn, error) {
	grpcDialOption := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	registryBuilder := clientcommon.NewGrpcBuilder(clusterClient.nodeData.RegistryServerAddress).
		SetGrpcDialOptions(grpcDialOption...)

	conn, err := registryBuilder.Build()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create registry connection: %w", err)
	}

	client := pb_registry.NewRegistryServiceClient(conn)
	return client, conn, nil
}

func retrieveAllNodesFromRegistry(cntNodeAddress, cntHostname string, registryClient pb_registry.RegistryServiceClient, clusterClient *ClusterClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nodeList, err := registryClient.GetNodeList(ctx, &pb_registry.NodeListRequest{})
	if err != nil {
		clusterClient.logger.Error("Error retrieving the node list from registry", "error", err)
		return fmt.Errorf("failed to get node list: %w", err)
	}

	// Lock only when modifying shared state
	clusterClient.nodeData.Mu.Lock()
	defer clusterClient.nodeData.Mu.Unlock()

	// Set up node list in map
	for _, node := range nodeList.NodeList {
		nodeName := node.NodeHostname
		nodeIP := node.NodeIP
		nodeControlPort := node.NodeControlPort

		// Skip self
		if nodeName == cntHostname && nodeIP == cntNodeAddress {
			continue
		}

		peerNode := nodecommon.InitializeNode(nodeName, nodeIP, nodeControlPort, "", 1)
		clusterClient.nodeData.PeerNodes[nodeName] = *peerNode

		clusterClient.logger.Info("Added peer node",
			"hostname", nodeName,
			"ip", nodeIP,
			"port", nodeControlPort)
	}

	clusterClient.logger.Info("Successfully retrieved peer nodes", "count", len(nodeList.NodeList)-1)
	return nil
}

func (clusterClient *ClusterClient) RegisterNodeWithRegistry(cntNodeAddress, cntHostname, cntControlPort, cntDataPort string) error {
	clusterClient.logger.Info("Starting node registration with registry")

	// Create client with proper resource management
	registryClient, conn, err := clusterClient.createRegistryClient()
	if err != nil {
		clusterClient.logger.Error("Error creating registry client", "error", err)
		return err
	}
	defer conn.Close() // Ensure connection is closed

	request := &pb_registry.RegisterNodeRequest{
		Hostname:   cntHostname,
		IpAddress:  cntNodeAddress,
		PortNumber: cntControlPort,
	}

	// Add timeout for registration
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Register the node
	response, err := registryClient.RegisterNode(ctx, request)
	if err != nil {
		clusterClient.logger.Error("Failed to register node with registry", "error", err)
		return fmt.Errorf("registration failed: %w", err)
	}

	clusterClient.logger.Info("Successfully registered with registry",
		"status", response.Status,
		"message", response.Message)

	// Handle error from node retrieval
	if err := retrieveAllNodesFromRegistry(cntNodeAddress, cntHostname, registryClient, clusterClient); err != nil {
		clusterClient.logger.Error("Failed to retrieve peer nodes", "error", err)
		return fmt.Errorf("failed to retrieve peer nodes: %w", err)
	}

	return nil
}

func (clusterClient *ClusterClient) SendRegularNodeHeartBeat(cntHostname, cntIP, cntCPPort string) {
	clusterClient.logger.Info("Starting heartbeat service")

	// Create client once and reuse
	registryClient, conn, err := clusterClient.createRegistryClient()
	if err != nil {
		clusterClient.logger.Error("Error creating registry client for heartbeat", "error", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	heartbeatRequest := &pb_registry.HeartBeatRequest{
		Hostname:   cntHostname,
		IpAddress:  cntIP,
		PortNumber: cntCPPort,
	}

	for range ticker.C {
		// Add timeout for each heartbeat
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		clusterClient.logger.Debug("Sending heartbeat to registry server")

		// Handle heartbeat response and errors
		response, err := registryClient.NodeHeartBeat(ctx, heartbeatRequest)
		if err != nil {
			clusterClient.logger.Error("Failed to send heartbeat", "error", err)
			cancel()
			continue
		}

		if response.Status != "200" {
			clusterClient.logger.Warn("Heartbeat returned non-success status",
				"status", response.Status,
				"message", response.Message)
		} else {
			clusterClient.logger.Debug("Successfully sent heartbeat")
		}

		cancel()
	}
}
