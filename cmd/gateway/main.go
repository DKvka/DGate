package main

import (
	"dgateway/internal/manager"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", ":8844", "port to serve")
	flag.Parse()

	logFile, err := os.Create("log.txt")
	if err != nil {
		log.Fatal("Error creating log file: ", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	m := manager.NewManager()

	go m.RunCli()

	mux := http.NewServeMux()

	mux.HandleFunc("/", m.Route)

	log.Fatal(http.ListenAndServe(*port, mux))
}
