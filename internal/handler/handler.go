package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/websocket"
)

// Creates a basic http handler that calls the
// destination backend with the incoming request
func Create(dest string) http.HandlerFunc {
	url, err := url.Parse(dest)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println()
		log.Println("Incoming request from:", r.RemoteAddr, " - Routing to:", url)
		proxy.ServeHTTP(w, r)
		log.Println("Roundtrip success, response sent to client")
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
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println()
		log.Println("Incoming request from:", r.RemoteAddr, " - Routing to:", url)
		proxy.ServeHTTP(w, r)
		log.Println("Roundtrip success, response sent to client")
	}
}
