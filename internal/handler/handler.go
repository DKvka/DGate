package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// Simpe ID system to track request logs
type idspawner struct {
	id uint64
	m  sync.Mutex
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
	_url, err := url.Parse(dest)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(_url)

	spawner := NewIdSpawner(1000000000000000)

	return func(w http.ResponseWriter, r *http.Request) {
		handleID := spawner.Next()

		log.Printf("%d - Incoming request from: %s - routing to: %s \n", handleID, r.RemoteAddr, _url)
		proxy.ServeHTTP(w, r)
		log.Printf("%d - Roundtrip success, response sent to client \n", handleID)
	}
}

// Creates an http handler that allows websocket upgrades
// to the destination backend
func CreateWithWebsocket(dest string) http.HandlerFunc {
	_url, err := url.Parse(dest)
	if err != nil {
		log.Fatal(err)
	}

	wsURL := *_url
	wsURL.Scheme = "ws"

	proxy := httputil.NewSingleHostReverseProxy(_url)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	spawner := NewIdSpawner(2000000000000000)

	return func(w http.ResponseWriter, r *http.Request) {
		handleID := spawner.Next()

		if strings.HasSuffix(r.URL.Path, "/ws") {
			// Upgrade client connection
			log.Printf("%d - Upgrading client to websocket...", handleID)
			clientConn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Printf("%d - %s", handleID, err)
				http.Error(w, "Error opening websocket", http.StatusInternalServerError)
				return
			}
			defer clientConn.Close()
			log.Printf("%d - Client successfully upgraded", handleID)

			// Connect websocket to backend
			log.Printf("%d - Connecting to backend server...", handleID)
			serverConn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), r.Header)
			if err != nil {
				log.Printf("%d - %s", handleID, err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			defer serverConn.Close()
			log.Printf("%d - Connection opened to backend", handleID)

			log.Printf("%d - Relaying messages between client and server", handleID)

			// Graceful shutdown channels for client and server connections
			killc := make(chan struct{})
			kills := make(chan struct{})

			// Client to server relay
			go func() {
				for {
					select {
					case <-killc:
						kills <- struct{}{}
						return
					default:
						t, msg, err := clientConn.ReadMessage()
						if err != nil {
							log.Printf("%d - %s", handleID, err)
							kills <- struct{}{}
						}

						err = serverConn.WriteMessage(t, msg)
						if err != nil {
							log.Printf("%d - %s", handleID, err)
							kills <- struct{}{}
						}
					}
				}
			}()

			// Server to client relay
			for {
				select {
				case <-kills:
					killc <- struct{}{}
					return
				default:
					t, msg, err := serverConn.ReadMessage()
					if err != nil {
						log.Printf("%d - %s", handleID, err)
						killc <- struct{}{}
					}

					err = clientConn.WriteMessage(t, msg)
					if err != nil {
						log.Printf("%d - %s", handleID, err)
						killc <- struct{}{}
					}
				}
			}
		}

		// Serve initial HTTP for the initial non websocket request
		log.Printf("%d - Incoming request from: %s - routing to: %s \n", handleID, r.RemoteAddr, _url)
		proxy.ServeHTTP(w, r)
		log.Printf("%d - Initial rountrip success, waiting for upgrade to websocket... \n", handleID)
	}
}
