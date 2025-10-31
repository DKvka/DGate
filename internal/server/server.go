package server

import (
	"dgateway/internal/config"
)

// Entry point
func Run(configPath string) error {
	cfg, err := config.Get(configPath)
	if err != nil {
		return err
	}

}
