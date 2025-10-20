package nodecommon

const (
	registry = iota
	workerNode
)

type Node struct {
	nodeHostname    string
	nodeIP          string
	nodeControlPort string
	nodeDataPort    string
	nodeType        int
}
