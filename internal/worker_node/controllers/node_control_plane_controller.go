package controllers

import (
	"fmt"
	"log/slog"

	nodecommon "github.com/Vahsek/distrokv/internal/common/node_common"
	"github.com/Vahsek/distrokv/internal/common/util"
	"github.com/Vahsek/distrokv/internal/worker_node/data"
	pb "github.com/Vahsek/distrokv/pkg/node/controlplane"
)

func RegisterNewPeerNode(request *pb.NewServerAddRequest, nodeData *data.NodeData, logger *slog.Logger) error {
	hostname := request.Hostname
	ipAddress := request.IpAddress
	controlPlanePort := request.ControlPlanePort
	dataPlanePort := request.DataPlanePort
	nodeHash := util.GenerateHash(hostname + ipAddress + controlPlanePort)
	logger.Info("Started registering the node: ", "NodeName: ", hostname, "NodeIP: ", ipAddress)
	nodeData.Mu.Lock()
	defer nodeData.Mu.Unlock()

	_, exists := nodeData.PeerNodes[nodeHash]

	if exists {
		logger.Error("Node Already Exists")
		return fmt.Errorf("Node Already Exists")
	}

	nodeData.PeerNodes[nodeHash] = *nodecommon.InitializeNode(
		hostname,
		ipAddress,
		controlPlanePort,
		dataPlanePort,
		1)

	logger.Info("Node Registered Successfully")
	return nil
}
