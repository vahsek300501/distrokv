package data

import (
	"log/slog"
	"sync"

	nodecommon "github.com/Vahsek/distrokv/internal/common/node_common"
)

type NodeData struct {
	NodeDetails           nodecommon.Node
	PeerNodes             map[string]nodecommon.Node
	RegistryServerAddress string
	Logger                slog.Logger
	Mu                    sync.Mutex
}
