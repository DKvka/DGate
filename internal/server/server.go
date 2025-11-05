package server

import (
	"dgateway/internal/config"
	"dgateway/internal/handler"
	"net/http"
)

// Entry point
func Run(cfg *config.Config) error {
	mux := http.NewServeMux()

	for _, servercfg := range cfg.ServerPool {
		if servercfg.AllowWebsocket {
			mux.Handle(servercfg.GatewayEndpoint, handler.CreateWithWebsocket(servercfg.Destination))
		} else {
			mux.Handle(servercfg.GatewayEndpoint, handler.Create(servercfg.Destination))
		}
	}

	return http.ListenAndServe(cfg.Gateway.Port, mux)
}
