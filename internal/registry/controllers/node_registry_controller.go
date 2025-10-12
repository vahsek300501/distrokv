package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
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

func RegisterNewNode(nodeDetails *pb.RegisterNodeRequest) bool {
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

	_, exists := RegisteredNodes.registeredNodesMap[nodeHash]

	if exists {
		return false
	}
	RegisteredNodes.registeredNodesMap[nodeHash] = newNode
	return true
}

func RegisterNodeHeartBeat(nodeDetails *pb.HeartBeatRequest) bool {
	RegisteredNodes.mu.Lock()
	defer RegisteredNodes.mu.Unlock()

	nodeHash := GetHash(nodeDetails.Hostname + nodeDetails.IpAddress)
	existingNode, exists := RegisteredNodes.registeredNodesMap[nodeHash]

	if !exists {
		return false
	}

	existingNode.lastHeartBeatTime = time.Now()
	RegisteredNodes.registeredNodesMap[nodeHash] = existingNode

	return true
}
