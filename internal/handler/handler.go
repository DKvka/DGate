package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
		proxy.ServeHTTP(w, r)
	}
}

// Creates an http handler that allows websocket upgrades
// to the destination backend
func CreateWithWebsocket(dest string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
