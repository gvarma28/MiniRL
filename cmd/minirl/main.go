package main

import (
	"net/http"
	"os"

	"github.com/gvarma28/MiniRL/internal/datastore"
	"github.com/gvarma28/MiniRL/internal/logger"
	"github.com/gvarma28/MiniRL/internal/proxy"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		BackendURL    string `yaml:"backend_url"`
		DatastoreType string `yaml:"datastore_type"`
	} `yaml:"server"`
	Redis struct {
		Address string `yaml:"address"`
	} `yaml:"redis"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	log := logger.NewLogger()
	defer log.Sync()

	config, err := LoadConfig("../../minirl-config.yaml")
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	datastoreType := config.Server.DatastoreType
	backendURL := config.Server.BackendURL
	if backendURL == "" {
		log.Fatal("BACKEND_URL not set")
	}

	store := datastore.NewDatastore(datastoreType, log)

	log.Info("Starting proxy server on :4000")
	if err := http.ListenAndServe(":4000", proxy.NewProxy(backendURL, store, log)); err != nil {
		log.Fatal("Server failed", zap.Error(err))
	}

}
