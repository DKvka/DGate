package config

import (
	"encoding/json"
	"io"
	"os"
)

type server struct {
	GatewayEndpoint string `json:"gateway_endpoint"`
	Name            string `json:"name"`
	Destination     string `json:"destination"`
}

type config struct {
	Gateway struct {
		Addr string `json:"addr"`
		Port string `json:"port"`
	}
	ServerPool []server
}

// Reads a JSON file into a `config` struct and returns it
func Get(filepath string) (*config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configjson, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg = new(config)
	err = json.Unmarshal(configjson, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
