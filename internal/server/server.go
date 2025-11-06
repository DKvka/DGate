package server

import (
	"dgate/internal/config"
	"dgate/internal/handler"
	"log"
	"net/http"
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
