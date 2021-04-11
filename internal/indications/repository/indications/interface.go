package indications

import (
	"net"

	"github.com/Systemnick/argo_exporter/internal/domain/indication"
)

type Indications *indication.Indication

type Repository interface {
	Poll(ip net.IP, metrics *MURMetrics) error
}

// type RepositoryFabric interface {
// 	New(ip net.IP, meterCount uint16)
// 	OpenConnections(onReady func())
// 	Poll(metrics *MURMetrics, onReady func())
// 	CloseConnections(onReady func())
// }
