package handler

import (
	"log"
	"math/rand"
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

func NewIdSpawner() *idspawner {
	return &idspawner{
		id: rand.Uint64(),
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

	spawner := NewIdSpawner()

	return func(w http.ResponseWriter, r *http.Request) {
		handleID := spawner.Next()

		log.Printf("%d - Incoming request from: %s - routing to: %s \n", handleID, r.RemoteAddr, _url)
		proxy.ServeHTTP(w, r)
		log.Printf("%d - Roundtrip success, response sent to client \n", handleID)
	}
}

// Creates a websocket handler
// Timeouts must be in seconds
func CreateWithWebsocket(dest string, buffersize, clientTimeout, serverTimeout int) http.HandlerFunc {
	_url, err := url.Parse(dest)
	if err != nil {
		log.Fatal(err)
	}

	wsURL := *_url
	wsURL.Scheme = "ws"
	if _url.Scheme == "https" {
		wsURL.Scheme = "wss"
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  buffersize,
		WriteBufferSize: buffersize,
		CheckOrigin: func(r *http.Request) bool {
			return true
			/*
				origin := r.Header.Get("Origin")
				if origin == "" {
					host := r.Host
					return host == "127.0.0.1" || host == "localhost"
				}

				u, err := url.Parse(origin)
				if err != nil {
					return false
				}

				return u.Hostname() == "127.0.0.1" ||
					u.Hostname() == "localhost" ||
					u.Hostname() == "deekays.ddns.net"
			*/
		},
	}

	spawner := NewIdSpawner()

	return func(w http.ResponseWriter, r *http.Request) {
		handleID := spawner.Next()

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
			wsURL.String(),
			header,
		)
		if err != nil {
			log.Printf("%d - %s", handleID, err)
			return
		}
		defer serverConn.Close()
		log.Printf("%d - Connection opened to backend", handleID)

		// Connection keep-alive/deadline setup
		serverConn.SetReadDeadline(time.Now().Add(time.Duration(serverTimeout) * time.Second))
		serverConn.SetPongHandler(func(string) error {
			serverConn.SetReadDeadline(time.Now().Add(time.Duration(serverTimeout) * time.Second))
			return nil
		})
		clientConn.SetReadDeadline(time.Now().Add(time.Duration(clientTimeout) * time.Second))
		clientConn.SetPongHandler(func(string) error {
			clientConn.SetReadDeadline(time.Now().Add(time.Duration(clientTimeout) * time.Second))
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
					t, msg, err := clientConn.ReadMessage()
					if err != nil {
						log.Printf("%d - %s", handleID, err)
						return
					}

					err = serverConn.WriteMessage(t, msg)
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
}
