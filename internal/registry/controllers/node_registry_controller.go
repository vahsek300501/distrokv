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

func (n node) GetHash() string {
	var plaintext string = n.hostname + n.ipAddress
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
	nodeHash := newNode.GetHash()

	RegisteredNodes.mu.Lock()
	defer RegisteredNodes.mu.Unlock()

	_, exists := RegisteredNodes.registeredNodesMap[nodeHash]

	if exists {
		return false
	}
	RegisteredNodes.registeredNodesMap[nodeHash] = newNode
	return true
}
