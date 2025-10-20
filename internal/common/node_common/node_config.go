package nodecommon

const (
	registry = iota
	workerNode
)

type Node struct {
	NodeHostname    string
	NodeIP          string
	NodeControlPort string
	NodeDataPort    string
	NodeType        int
}

func InitializeNode(hostname, IP, nodeCPNumber, nodeDPNumber string, nodeType int) *Node {
	return &Node{
		NodeHostname:    hostname,
		NodeIP:          IP,
		NodeControlPort: nodeCPNumber,
		NodeDataPort:    nodeDPNumber,
		NodeType:        nodeType,
	}
}
