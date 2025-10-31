package manager

import (
	"net/http/httputil"
	"sync"
)

// Manager manages the gateway server, routes requests and
// runs a CLI for manual runtime modifications to the gateway
type Manager struct {
	// maps a service name to its address
	services map[string]*httputil.ReverseProxy
	lock     sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		services: make(map[string]*httputil.ReverseProxy),
	}
}

func (m *Manager) AddService(name, addr string) {
	m.lock.Lock()
	defer m.lock.Unlock()
}

func (m *Manager) DeleteService(name string) {
	m.lock.Lock()
	defer m.lock.Unlock()
}

