package common

import (
	"AdvanceTradeEngine/models"
	"log"

	"github.com/BurntSushi/toml"
)

const (
	Error   = "E"
	Success = "S"
)

func LoadConfig(path string) *models.Config {
	var cfg models.Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}
	return &cfg
}
