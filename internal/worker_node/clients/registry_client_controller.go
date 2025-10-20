package clients

import (
	"context"

	clientcommon "github.com/Vahsek/distrokv/internal/common/client_common"
	pb_registry "github.com/Vahsek/distrokv/pkg/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func retriveAllNodesFromRegistry(cntNodeAddress, cntHostname string, registryClient pb_registry.RegistryServiceClient, clusterClient *ClusterClient) error {
	nodeList, err := registryClient.GetNodeList(context.Background(), &pb_registry.NodeListRequest{})
	if err != nil {
		clusterClient.logger.Error("Error retrving the node list from registry", "error", err)
		return err
	}

	// Set up node list in map
	for _, node := range nodeList.NodeList {
		node_name := node.NodeHostname
		node_ip := node.NodeIP
		node_control_port := node.NodeControlPort
		if node_name == cntHostname && node_ip == cntNodeAddress {
			continue
		}
		peerNode := PeerNodes{
			hostname:    node_name,
			ipAddress:   node_ip,
			controlPort: node_control_port,
		}
		clusterClient.nodes[node_name] = peerNode
	}

	return nil
}

func (clusterClient *ClusterClient) RegisterNodeWithRegistry(cntNodeAddress, cntHostname, cntDataPort, cntControlPort string) error {
	clusterClient.mu.Lock()
	defer clusterClient.mu.Unlock()

	// Creating connection to Registry Server for registering the nodes
	grpcDialOption := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	registryBuilder := clientcommon.NewGrpcBuilder(clusterClient.registryAddress).SetGrpcDialOptions(grpcDialOption...)
	registryClientConstructor := func(conn *grpc.ClientConn) pb_registry.RegistryServiceClient {
		return pb_registry.NewRegistryServiceClient(conn)
	}
	registryClient, err := clientcommon.GetClient(context.Background(), clusterClient.factory, *registryBuilder, registryClientConstructor)
	if err != nil {
		clusterClient.logger.Error("Error creating registry client")
		return err
	}

	request := &pb_registry.RegisterNodeRequest{
		Hostname:   cntHostname,
		IpAddress:  cntNodeAddress,
		PortNumber: cntControlPort,
	}

	// Register the node
	response, err := registryClient.RegisterNode(context.Background(), request)
	if err != nil {
		clusterClient.logger.Error("Failed to register node with registry", "error", err)
		return err
	}
	clusterClient.logger.Info("Successfully registered with registry", "status", response.Status, "message", response.Message)

	// Get All the nodes from the registry
	retriveAllNodesFromRegistry(cntNodeAddress, cntHostname, registryClient, clusterClient)
	return nil
}
