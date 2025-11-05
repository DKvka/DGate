package main

import (
	"dgateway/internal/config"
	"dgateway/internal/server"
	"flag"
	"log"
)

func main() {
	configPath := flag.String("cpath", "config.json", "sets path to configuration json file")
	flag.Parse()

	cfg, err := config.Get(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Run(cfg))
}
