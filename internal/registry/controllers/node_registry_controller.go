package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log/slog"
	"sync"
	"time"

	pb "github.com/Vahsek/distrokv/pkg/registry"
)

type node struct {
	hostname          string
	ipAddress         string
	portNumber        string
	registrationTime  time.Time
	lastHeartBeatTime time.Time
}

func GetHash(plaintext string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, plaintext)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

type NodeMap struct {
	registeredNodesMap map[string]node
	mu                 sync.Mutex
}

var RegisteredNodes NodeMap = NodeMap{
	registeredNodesMap: make(map[string]node),
}

func RegisterNewNode(nodeDetails *pb.RegisterNodeRequest, logger slog.Logger) bool {
	logger.Info("Creating new Node with Nodename: %s NodeIP: %s", nodeDetails.Hostname, nodeDetails.IpAddress)
	var newNode node = node{
		hostname:          nodeDetails.Hostname,
		ipAddress:         nodeDetails.IpAddress,
		portNumber:        nodeDetails.PortNumber,
		registrationTime:  time.Now(),
		lastHeartBeatTime: time.Now(),
	}
	nodeHash := GetHash(newNode.hostname + newNode.ipAddress)

	RegisteredNodes.mu.Lock()
	defer RegisteredNodes.mu.Unlock()

	logger.Info("Checking the node hash")
	_, exists := RegisteredNodes.registeredNodesMap[nodeHash]

	if exists {
		logger.Error("Node hash already exists. Skipping adding the node")
		return false
	}
	logger.Info("Adding node to node dictionary")
	RegisteredNodes.registeredNodesMap[nodeHash] = newNode
	return true
}

func RegisterNodeHeartBeat(nodeDetails *pb.HeartBeatRequest, logger slog.Logger) bool {
	RegisteredNodes.mu.Lock()
	defer RegisteredNodes.mu.Unlock()

	logger.Info("Heatbeat for node with Nodename: %s NodeIP: %s", nodeDetails.Hostname, nodeDetails.IpAddress)
	nodeHash := GetHash(nodeDetails.Hostname + nodeDetails.IpAddress)
	existingNode, exists := RegisteredNodes.registeredNodesMap[nodeHash]

	if !exists {
		logger.Error("Node Doesn't exists")
		return false
	}
	logger.Info("Registering node heartbeat")
	existingNode.lastHeartBeatTime = time.Now()
	RegisteredNodes.registeredNodesMap[nodeHash] = existingNode

	return true
}
