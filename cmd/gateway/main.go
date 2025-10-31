package main

import (
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

	mux := http.NewServeMux()

	mux.HandleFunc("/", nil)

	log.Fatal(http.ListenAndServe(*port, mux))
}

