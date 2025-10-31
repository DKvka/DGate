package main

import (
	"dgateway/internal/server"
	"flag"
	"log"
)

func main() {
	configpath := flag.String("cpath", "config.json", "sets path to configuration json file")
	flag.Parse()

	log.Fatal(server.Run(*configpath))
}
