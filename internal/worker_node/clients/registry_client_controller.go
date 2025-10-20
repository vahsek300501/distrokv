package clients

import (
	"context"
	"fmt"
	"time"

	nodecommon "github.com/Vahsek/distrokv/internal/common/node_common"
	"github.com/Vahsek/distrokv/internal/worker_node/data"
	pb_registry "github.com/Vahsek/distrokv/pkg/registry"
)

func retrieveAllNodesFromRegistry(nodeData *data.NodeData, registryClient pb_registry.RegistryServiceClient, clusterClient *ClusterClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nodeList, err := registryClient.GetNodeList(ctx, &pb_registry.NodeListRequest{})
	if err != nil {
		clusterClient.logger.Error("Error retrieving the node list from registry", "error", err)
		return fmt.Errorf("failed to get node list: %w", err)
	}

	// Lock only when modifying shared state
	nodeData.Mu.Lock()
	defer nodeData.Mu.Unlock()

	// Set up node list in map
	for _, node := range nodeList.NodeList {
		nodeName := node.NodeHostname
		nodeIP := node.NodeIP
		nodeControlPort := node.NodeControlPort

		// Skip self
		if nodeName == nodeData.NodeDetails.NodeHostname && nodeIP == nodeData.NodeDetails.NodeIP {
			continue
		}

		peerNode := nodecommon.InitializeNode(nodeName, nodeIP, nodeControlPort, "", 1)
		nodeData.PeerNodes[nodeName] = *peerNode

		clusterClient.logger.Info("Added peer node",
			"hostname", nodeName,
			"ip", nodeIP,
			"port", nodeControlPort)
	}

	clusterClient.logger.Info("Successfully retrieved peer nodes", "count", len(nodeList.NodeList)-1)
	return nil
}

func (clusterClient *ClusterClient) RegisterNodeWithRegistry(nodeData *data.NodeData) error {
	clusterClient.logger.Info("Starting node registration with registry")

	// Create client with proper resource management
	registryClient, err := clusterClient.createRegistryClient(nodeData)
	if err != nil {
		clusterClient.logger.Error("Error creating registry client", "error", err)
		return err
	}

	request := &pb_registry.RegisterNodeRequest{
		Hostname:   nodeData.NodeDetails.NodeHostname,
		IpAddress:  nodeData.NodeDetails.NodeIP,
		PortNumber: nodeData.NodeDetails.NodeControlPort,
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
	if err := retrieveAllNodesFromRegistry(nodeData, registryClient, clusterClient); err != nil {
		clusterClient.logger.Error("Failed to retrieve peer nodes", "error", err)
		return fmt.Errorf("failed to retrieve peer nodes: %w", err)
	}

	return nil
}

func (clusterClient *ClusterClient) SendRegularNodeHeartBeat(nodeData *data.NodeData) {
	clusterClient.logger.Info("Starting heartbeat service")

	// Create client once and reuse
	registryClient, err := clusterClient.createRegistryClient(nodeData)
	if err != nil {
		clusterClient.logger.Error("Error creating registry client for heartbeat", "error", err)
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	heartbeatRequest := &pb_registry.HeartBeatRequest{
		Hostname:   nodeData.NodeDetails.NodeHostname,
		IpAddress:  nodeData.NodeDetails.NodeIP,
		PortNumber: nodeData.NodeDetails.NodeControlPort,
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
