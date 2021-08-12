package config

import (
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	configFile = "../config.json"
	InitConfig()
	log.Println("test")
}
