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
		log.Printf("Adding server to pool: %s - endpoint at gateway: %s", sCfg.Name, sCfg.GatewayEndpoint)
		if sCfg.AllowWebsocket {
			mux.HandleFunc(sCfg.GatewayEndpoint, handler.Create(sCfg.Destination))
			mux.HandleFunc(sCfg.GatewayEndpoint+sCfg.WebsockSettings.WebsockEndpoint,
				handler.CreateWithWebsocket(
					sCfg.Destination+sCfg.WebsockSettings.WebsockEndpoint,
					2048,
					25,
					25,
				))
		} else {
			mux.HandleFunc(sCfg.GatewayEndpoint, handler.Create(sCfg.Destination))
		}
	}

	log.Println()
	log.Println("Starting up server on port ", cfg.Gateway.Port)
	return http.ListenAndServe(cfg.Gateway.Port, mux)
}
