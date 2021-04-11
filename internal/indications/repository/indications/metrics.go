package indications

import (
	"net"
	"sync"
	"time"
)

type MURMetrics struct {
	RegStatus   map[string]int
	RegDuration map[string]time.Duration
	Meters      map[string]map[string]string
	mu          sync.Mutex
}

func NewMURMetrics() *MURMetrics {
	return &MURMetrics{
		RegStatus:   make(map[string]int),
		RegDuration: make(map[string]time.Duration),
		Meters:      make(map[string]map[string]string),
	}
}

func (m *MURMetrics) SetRegStatus(ip net.IP, i int) {
	m.mu.Lock()
	m.RegStatus[ip.String()] = i
	m.mu.Unlock()
}

func (m *MURMetrics) SetRegPollDuration(ip net.IP, d time.Duration) {
	m.mu.Lock()
	m.RegDuration[ip.String()] = d
	m.mu.Unlock()
}

func (m *MURMetrics) SetMetersIndications(in map[string]map[string]string) {
	m.mu.Lock()
	m.Meters = mergeMaps(m.Meters, in)
	m.mu.Unlock()
}

func mergeMaps(ms ...map[string]map[string]string) map[string]map[string]string {
	res := map[string]map[string]string{}

	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}

	return res
}
