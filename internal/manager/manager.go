package manager

import (
	"net/http/httputil"
	"sync"
)

// Manager manages the gateway server, routes requests and
// runs a CLI for manual modifications to the gateway
type Manager struct {
	// maps a service name to its address
	services map[string]string
	proxy    httputil.ReverseProxy
	lock     sync.Mutex
}

func NewManager() Manager {
	return Manager{
		services: make(map[string]string),
	}
}

func (m *Manager) AddService() {
	m.lock.Lock()
	defer m.lock.Unlock()
}

func (m *Manager) DeleteService() {
	m.lock.Lock()
	defer m.lock.Unlock()
}
