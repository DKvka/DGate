package server

import (
	"dgateway/internal/config"
	"dgateway/internal/handler"
	"net/http"
	"log"
)

// Entry point
func Run(cfg *config.Config) error {
	mux := http.NewServeMux()

	log.Println()
	log.Println("Readying resources...")
	for _, servercfg := range cfg.ServerPool {
		log.Println("Adding server to pool:", servercfg.Name)
		if servercfg.AllowWebsocket {
			mux.Handle(servercfg.GatewayEndpoint, handler.CreateWithWebsocket(servercfg.Destination))
		} else {
			mux.Handle(servercfg.GatewayEndpoint, handler.Create(servercfg.Destination))
		}
	}

	log.Println()
	log.Println("Starting up server on port ", cfg.Gateway.Port)
	return http.ListenAndServe(cfg.Gateway.Port, mux)
}
