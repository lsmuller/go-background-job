package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type MainConfig struct {
	App    *AppConfig    `yaml:"app"`
	Worker *WorkerConfig `yaml:"worker"`
}

type AppConfig struct {
	Stack   string         `yaml:"stack"`
	Logging *LoggingConfig `yaml:"logging"`
}

type LoggingConfig struct {
	Verbose int  `yaml:"verbose"`
	LogJSON bool `yaml:"logJSON"`
}

type RedisConfig struct {
	Url  string `yaml:"url"`
	Pool string `yaml:"pool"`
}

type RetryConfig struct {
	Enabled  bool `yaml:"enabled"`
	MaxDelay int  `yaml:"maxDelay"`
	MinDelay int  `yaml:"minDelay"`
	RetryMax int  `yaml:"retryMax"`
	Exp      int  `yaml:"exp"`
	MaxRand  int  `yaml:"maxRand"`
}

type JobConfig struct {
	Enabled     bool              `yaml:"enabled"`
	Concurrency int               `yaml:"concurrency"`
	Period      time.Duration     `yaml:"period"`
	Metadata    map[string]string `yaml:"metadata"`
}

type JobsConfig struct {
	BitcoinPriceIndexFetcher *JobConfig `yaml:"bitcoinPriceIndexFetcher"`
}

type WorkerConfig struct {
	Retry *RetryConfig `yaml:"retry"`
	Redis *RedisConfig `yaml:"redis"`
	Jobs  *JobsConfig  `yaml:"jobs"`
}

func GetConfig() *MainConfig {
	yamlFile, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		panic("error reading yamlFile ")
	}
	var cfg *MainConfig

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic("error reading yamlFile ")
	}

	return cfg
}
