package controllers

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	nodecommon "github.com/Vahsek/distrokv/internal/common/node_common"
	"github.com/Vahsek/distrokv/internal/common/util"
	pb "github.com/Vahsek/distrokv/pkg/registry"
)

type RegisteredNodeDetails struct {
	nodeDetails       nodecommon.Node
	registrationTime  time.Time
	lastHeartBeatTime time.Time
}

type NodeRegistryInterface interface {
	RegisterNode(nodeDetails *pb.RegisterNodeRequest) error
	UpdateHeartbeat(nodeDetails *pb.HeartBeatRequest) error
	GetNodeList() []*pb.NodeDetails
}

type NodeRegistry struct {
	nodes  map[string]RegisteredNodeDetails
	mu     sync.Mutex
	logger slog.Logger
}

func InitializeNodeRegistry(logger slog.Logger) *NodeRegistry {
	return &NodeRegistry{
		nodes:  make(map[string]RegisteredNodeDetails),
		logger: logger,
	}
}

func (nodeRegistry *NodeRegistry) RegisterNewNode(nodeDetails *pb.RegisterNodeRequest) error {
	nodeRegistry.mu.Lock()
	defer nodeRegistry.mu.Unlock()

	nodeRegistry.logger.Info("Creating new Node with Nodename: %s NodeIP: %s", nodeDetails.Hostname, nodeDetails.IpAddress)
	var newNodeDetails nodecommon.Node = *nodecommon.InitializeNode(nodeDetails.Hostname, nodeDetails.IpAddress, nodeDetails.PortNumber, "8080", 1)
	var newNode RegisteredNodeDetails = RegisteredNodeDetails{
		nodeDetails:       newNodeDetails,
		registrationTime:  time.Now(),
		lastHeartBeatTime: time.Now(),
	}
	nodeHash := util.GenerateHash(newNode.nodeDetails.NodeHostname + newNode.nodeDetails.NodeIP)

	nodeRegistry.logger.Info("Checking the node hash")
	_, exists := nodeRegistry.nodes[nodeHash]

	if exists {
		nodeRegistry.logger.Error("Node hash already exists. Skipping adding the node")
		return fmt.Errorf("Node hash already exists. Skipping adding the node")
	}
	nodeRegistry.logger.Info("Adding node to node dictionary")
	nodeRegistry.nodes[nodeHash] = newNode
	return nil
}

func (nodeRegistry *NodeRegistry) RegisterNodeHeartBeat(nodeDetails *pb.HeartBeatRequest) error {
	nodeRegistry.mu.Lock()
	defer nodeRegistry.mu.Unlock()

	nodeRegistry.logger.Info("Heatbeat for node with Nodename: %s NodeIP: %s", nodeDetails.Hostname, nodeDetails.IpAddress)
	nodeHash := util.GenerateHash(nodeDetails.Hostname + nodeDetails.IpAddress)
	existingNode, exists := nodeRegistry.nodes[nodeHash]

	if !exists {
		nodeRegistry.logger.Error("Node Doesn't exists")
		return fmt.Errorf("Node Doesn't exists")
	}
	nodeRegistry.logger.Info("Registering node heartbeat")
	existingNode.lastHeartBeatTime = time.Now()
	nodeRegistry.nodes[nodeHash] = existingNode
	return nil
}

func (nodeRegistry *NodeRegistry) GetNodeList() []*pb.NodeDetails {
	nodeRegistry.mu.Lock()
	defer nodeRegistry.mu.Unlock()

	nodeRegistry.logger.Info("Generating Node List response")

	var nodes []*pb.NodeDetails
	registeredNodesMap := nodeRegistry.nodes
	for _, value := range registeredNodesMap {
		nodeDetail := &pb.NodeDetails{
			NodeIP:          value.nodeDetails.NodeIP,
			NodeHostname:    value.nodeDetails.NodeHostname,
			NodeControlPort: value.nodeDetails.NodeControlPort,
		}
		nodes = append(nodes, nodeDetail)
	}
	return nodes
}
