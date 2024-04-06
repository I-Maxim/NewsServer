package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Config struct {
	PollPeriod     time.Duration `yaml:"pollPeriod"`
	PollBatchSize  uint          `yaml:"pollBatchSize"`
	MongoURI       string        `yaml:"mongoURI"`
	DbName         string        `yaml:"dbName"`
	CollectionName string        `yaml:"collectionName"`
	ServerAddr     string        `yaml:"serverAddr"`
}

func Load() (*Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "config/config.yaml", "path to config yaml file")
	flag.Parse()

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	cfg := Config{}
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file failed: %w", err)
	}

	return &cfg, nil
}
