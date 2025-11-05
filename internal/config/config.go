package config

import (
	"encoding/json"
	"io"
	"os"
)

// Stores information about a single server behind the gateway
type server struct {
	GatewayEndpoint string `json:"gateway_endpoint"`
	Name            string `json:"name"`
	Destination     string `json:"destination"`
	AllowWebsocket       bool `json:"allow_websocket"`
}

// Stores the gateway configuration info
type Config struct {
	Gateway struct {
		Addr string `json:"addr"`
		Port string `json:"port"`
	}
	ServerPool []server
}

// Reads a JSON file into a `config` struct and returns it
func Get(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configjson, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	err = json.Unmarshal(configjson, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
