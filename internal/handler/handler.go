package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// Simpe ID system to track request logs
type idspawner struct {
	id uint64
	m sync.Mutex
}

func NewIdSpawner(startvalue uint64) *idspawner {
	return &idspawner{
		id: startvalue,
	}
}

func (s *idspawner) Next() uint64 {
	s.m.Lock()
	defer s.m.Unlock()

	s.id++
	id := s.id

	return id
}

// Creates a basic http handler that calls the
// destination backend with the incoming request
func Create(dest string) http.HandlerFunc {
	url, err := url.Parse(dest)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	spawner := NewIdSpawner(1000000000000000)

	return func(w http.ResponseWriter, r *http.Request) {
		handleID := spawner.Next()

		log.Println()
		log.Println(handleID, " - ", "Incoming request from:", r.RemoteAddr, " - Routing to:", url)
		proxy.ServeHTTP(w, r)
		log.Println(handleID, " - ", "Roundtrip success, response sent to client")
	}
}

// Creates an http handler that allows websocket upgrades
// to the destination backend
func CreateWithWebsocket(dest string) http.HandlerFunc {
	url, err := url.Parse(dest)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == "http://127.0.0.1"
		},
	}

	spawner := NewIdSpawner(2000000000000000)

	return func(w http.ResponseWriter, r *http.Request) {
		handleID := spawner.Next()

		log.Println()
		log.Println(handleID, " - ", "Incoming request from:", r.RemoteAddr, " - Routing to:", url)
		proxy.ServeHTTP(w, r)
		log.Println(handleID, " - ", "Initial rountrip success, waiting for upgrade to websocket...")
	}
}
