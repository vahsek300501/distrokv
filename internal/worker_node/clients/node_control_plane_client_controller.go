package clients

import (
	"context"
	"fmt"

	"github.com/Vahsek/distrokv/internal/worker_node/data"
	pb_contol_plane "github.com/Vahsek/distrokv/pkg/node/controlplane"
)

func (clusterClient *ClusterClient) RegisterNodeWithPeers(nodeData *data.NodeData) error {
	nodeData.Mu.RLock()
	defer nodeData.Mu.RUnlock()
	logger := clusterClient.logger
	logger.Info("Sending alive message to peers")
	for key, value := range nodeData.PeerNodes {
		logger.Info("Sending alive message to Peer: ", "PeerHostname", key)
		request := &pb_contol_plane.NewServerAddRequest{
			Hostname:         nodeData.NodeDetails.NodeHostname,
			IpAddress:        nodeData.NodeDetails.NodeIP,
			ControlPlanePort: nodeData.NodeDetails.NodeControlPort,
			DataPlanePort:    nodeData.NodeDetails.NodeDataPort,
		}
		peerAddress := value.NodeIP + ":" + value.NodeControlPort

		peerNodeClient, err := clusterClient.createPeerClientConnection(peerAddress)
		if err != nil {
			logger.Error("Error in creating peer client")
			return fmt.Errorf("Error in creating peer client")
		}

		response, peerError := peerNodeClient.RegisterNewPeerServer(context.Background(), request)
		logger.Info("Peer Response", "Response", response)
		if peerError != nil {
			logger.Error("Network Error in registering the peer client")
			return peerError
		}

		logger.Info("Successfully registered peer client")
	}
	return nil
}
