package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

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
//
// Timeouts must be in seconds
func CreateWithWebsocket(dest string, buffersize, clientTimeout, serverTimeout int) http.HandlerFunc {
	_url, err := url.Parse(dest)
	if err != nil {
		log.Fatal(err)
	}

	wsURL := *_url
	wsURL.Scheme = "wss"
	if _url.Scheme == "http" {
		wsURL.Scheme = "ws"
	}

	proxy := httputil.NewSingleHostReverseProxy(_url)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  buffersize,
		WriteBufferSize: buffersize,
		CheckOrigin: func(r *http.Request) bool {
			u, err := url.Parse(r.Header.Get("Origin"))
			if err != nil {
				return false
			}

			return u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost" || u.Hostname() == "deekays.ddns.net"
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
				return
			}
			defer clientConn.Close()
			log.Printf("%d - Client successfully upgraded", handleID)

			// Connect websocket to backend with stripped header
			log.Printf("%d - Connecting to backend server...", handleID)

			header := http.Header{}
			for k, v := range r.Header {
				if strings.HasPrefix(strings.ToLower(k), "sec-websocket") ||
					strings.ToLower(k) == "connection" ||
					strings.ToLower(k) == "upgrade" {
					continue
				}
				header[k] = v
			}

			serverConn, _, err := websocket.DefaultDialer.Dial(
				wsURL.ResolveReference(r.URL).String(),
				header,
			)
			if err != nil {
				log.Printf("%d - %s", handleID, err)
				return
			}
			defer serverConn.Close()
			log.Printf("%d - Connection opened to backend", handleID)

			// Connection deadline setup
			serverConn.SetReadDeadline(time.Now().Add(time.Duration(serverTimeout) * time.Minute))
			serverConn.SetPongHandler(func(string) error {
				serverConn.SetReadDeadline(time.Now().Add(time.Duration(serverTimeout) * time.Minute))
				return nil
			})
			clientConn.SetReadDeadline(time.Now().Add(time.Duration(clientTimeout) * time.Minute))
			clientConn.SetPongHandler(func(string) error {
				clientConn.SetReadDeadline(time.Now().Add(time.Duration(clientTimeout) * time.Minute))
				return nil
			})

			log.Printf("%d - Relaying messages between client and server", handleID)

			// Graceful shutdown channel
			kill := make(chan struct{})
			var once sync.Once
			STOP := func() { once.Do(func() { close(kill) }) }

			// Client to server relay
			go func() {
				defer STOP()
				for {
					select {
					case <-kill:
						return
					default:
						t, msg, err := serverConn.ReadMessage()
						if err != nil {
							log.Printf("%d - %s", handleID, err)
							return
						}

						err = clientConn.WriteMessage(t, msg)
						if err != nil {
							log.Printf("%d - %s", handleID, err)
							return
						}
					}
				}
			}()

			// Server to client relay
			defer STOP()
			for {
				select {
				case <-kill:
					return
				default:
					t, msg, err := serverConn.ReadMessage()
					if err != nil {
						log.Printf("%d - %s", handleID, err)
						return
					}

					err = clientConn.WriteMessage(t, msg)
					if err != nil {
						log.Printf("%d - %s", handleID, err)
						return
					}
				}
			}
		}

		// Serve HTTP for the initial non websocket request
		log.Printf("%d - Incoming request from: %s - routing to: %s \n", handleID, r.RemoteAddr, _url)
		proxy.ServeHTTP(w, r)
		log.Printf("%d - Initial rountrip success, waiting for upgrade to websocket... \n", handleID)
	}
}
