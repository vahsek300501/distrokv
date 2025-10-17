package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"sync"
	"time"

	pb "github.com/Vahsek/distrokv/pkg/registry"
)

type Node struct {
	hostname          string
	ipAddress         string
	portNumber        string
	registrationTime  time.Time
	lastHeartBeatTime time.Time
}

type NodeRegistryInterface interface {
	RegisterNode(nodeDetails *pb.RegisterNodeRequest) error
	UpdateHeartbeat(nodeDetails *pb.HeartBeatRequest) error
	GetNodeList() []*pb.NodeDetails
}

type NodeRegistry struct {
	nodes  map[string]Node
	mu     sync.Mutex
	logger slog.Logger
}

func InitializeNodeRegistry(logger slog.Logger) *NodeRegistry {
	return &NodeRegistry{
		nodes:  make(map[string]Node),
		logger: logger,
	}
}

func (nodeRegistry *NodeRegistry) generateHash(hostname, ipAddress string) string {
	plaintext := hostname + ipAddress
	hasher := sha256.New()
	hasher.Write([]byte(plaintext))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (nodeRegistry *NodeRegistry) RegisterNewNode(nodeDetails *pb.RegisterNodeRequest) error {
	nodeRegistry.mu.Lock()
	defer nodeRegistry.mu.Unlock()

	nodeRegistry.logger.Info("Creating new Node with Nodename: %s NodeIP: %s", nodeDetails.Hostname, nodeDetails.IpAddress)
	var newNode Node = Node{
		hostname:          nodeDetails.Hostname,
		ipAddress:         nodeDetails.IpAddress,
		portNumber:        nodeDetails.PortNumber,
		registrationTime:  time.Now(),
		lastHeartBeatTime: time.Now(),
	}
	nodeHash := nodeRegistry.generateHash(newNode.hostname, newNode.ipAddress)

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
	nodeHash := nodeRegistry.generateHash(nodeDetails.Hostname, nodeDetails.IpAddress)
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
			NodeIP:          value.ipAddress,
			NodeHostname:    value.hostname,
			NodeControlPort: value.portNumber,
		}
		nodes = append(nodes, nodeDetail)
	}
	return nodes
}
