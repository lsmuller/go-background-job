package test

import (
	"github.com/lsmuller/go-background-job/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func GetConfig() *config.MainConfig {
	yamlFile, err := ioutil.ReadFile("./../test/config.yaml")
	if err != nil {
		panic("error reading yamlFile ")
	}
	var cfg *config.MainConfig

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic("error reading yamlFile ")
	}

	return cfg
}
