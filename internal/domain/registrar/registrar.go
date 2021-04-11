package registrar

import (
	"net"
	"sync"
)

type Registrar struct {
	IP           net.IP
	Vendor       string
	Model        string
	MeterCount   uint16
	isPollingNow bool
	mu           sync.Mutex
}

func (r *Registrar) IsPollingNow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.isPollingNow
}
