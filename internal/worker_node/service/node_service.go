package service

import (
	"log/slog"
	"sync"

	nodecommon "github.com/Vahsek/distrokv/internal/common/node_common"
	"github.com/Vahsek/distrokv/internal/storage"
	"github.com/Vahsek/distrokv/internal/worker_node/clients"
	"github.com/Vahsek/distrokv/internal/worker_node/data"
	"github.com/Vahsek/distrokv/internal/worker_node/servers"
)

type WorkerNodeService struct {
	NodeConfig      *nodecommon.Node
	NodeData        *data.NodeData
	ClusterClient   *clients.ClusterClient
	Storage         *storage.KeyValueStore
	RegistryAddress string
	logger          slog.Logger
}

func InitializeNewNodeService(hostname, ip, controlPort, dataPort string, nodeType int, registryAddress string, logger slog.Logger) *WorkerNodeService {
	nodeConfig := nodecommon.InitializeNode(hostname, ip, controlPort, dataPort, nodeType)
	nodeData := &data.NodeData{
		NodeDetails:           *nodeConfig,
		PeerNodes:             make(map[string]nodecommon.Node),
		RegistryServerAddress: registryAddress,
		Logger:                logger,
	}
	return &WorkerNodeService{
		NodeConfig:      nodeConfig,
		NodeData:        nodeData,
		ClusterClient:   clients.InitializeClusterClient(logger),
		Storage:         storage.NewKeyValueStore(logger),
		RegistryAddress: registryAddress,
		logger:          logger,
	}
}

func (nodeService *WorkerNodeService) BootstrapWorkerNode() {
	var wg sync.WaitGroup
	nodeService.logger.Info("Bootstrapping the worker node")
	nodeService.logger.Info("Registering with the registry server")
	err := nodeService.ClusterClient.RegisterNodeWithRegistry(
		nodeService.NodeData)
	if err != nil {
		nodeService.logger.Error("Error in registering node with the registry")
		return
	}
	nodeService.logger.Info("Successfully registered with registry")

	nodeService.logger.Info("Starting Control Plane Server")
	wg.Add(1)
	go func() {
		defer wg.Done()
		servers.StartNodeControlPlaneServer(
			":"+nodeService.NodeConfig.NodeControlPort,
			nodeService.logger,
			nodeService.ClusterClient,
			nodeService.NodeData)
	}()

	nodeService.logger.Info("Starting Data Plane Server")
	wg.Add(1)
	go func() {
		defer wg.Done()
		servers.StartNodeDataPlaneServer(
			":"+nodeService.NodeConfig.NodeDataPort,
			nodeService.logger)
	}()

	nodeService.logger.Info("Starting heartbeat service")
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				nodeService.logger.Error("Heartbeat service panicked", "error", r)
			}
		}()
		nodeService.ClusterClient.SendRegularNodeHeartBeat(
			nodeService.NodeData)
	}()

	nodeService.logger.Info("Registering Node with peers")
	errRegisteringWithPeers := nodeService.ClusterClient.RegisterNodeWithPeers(nodeService.NodeData)
	if errRegisteringWithPeers != nil {
		nodeService.logger.Error("failed to register service with peers")
		return
	}
	nodeService.logger.Info("Successfully registerd with peers")
	wg.Wait()
}
