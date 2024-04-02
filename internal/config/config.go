package config

import "time"

type Config struct {
	PollPeriod     time.Duration `yaml:"pollPeriod"`
	PollBatchSize  uint          `yaml:"pollBatchSize"`
	MongoURI       string        `yaml:"mongoURI"`
	DbName         string        `yaml:"dbName"`
	CollectionName string        `yaml:"collectionName"`
	ServerAddr     string        `yaml:"serverAddr"`
}
