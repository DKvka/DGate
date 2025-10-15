package manager

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (m *Manager) Route(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprintln(w, "Root request received")
		return
	}

	target := strings.Split(r.URL.Path, "/")[0]

	m.lock.Lock()
	service := m.services[target]
	m.lock.Unlock()

	serviceURL, err := url.Parse(service)
	if err != nil {
		http.Error(w, "Requested service not found", http.StatusNotFound)
		return
	}
	r.URL = serviceURL

	resp, err := m.proxy.Transport.RoundTrip(r)
	if err != nil {

	}
}
