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
	for _, sCfg := range cfg.ServerPool {
		log.Println("Adding server to pool:", sCfg.Name)
		if sCfg.AllowWebsocket {
			mux.Handle(sCfg.GatewayEndpoint, handler.CreateWithWebsocket(
				sCfg.Destination,
				2048,
				sCfg.WebsockSettings.ServerTimeout,
				sCfg.WebsockSettings.ClientTimeout))
		} else {
			mux.Handle(sCfg.GatewayEndpoint, handler.Create(sCfg.Destination))
		}
	}

	log.Println()
	log.Println("Starting up server on port ", cfg.Gateway.Port)
	return http.ListenAndServe(cfg.Gateway.Port, mux)
}
