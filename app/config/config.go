package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"
)

type Config struct {
	AppPort string `yaml:"appPort"`
	DatabaseUser string `yaml:"databaseUser"`
	DatabaseName string `yaml:"databaseName"`
}

// GetConfig generates settings from config file
func GetConfig(path string) (Config, error) {
	log.Println(path)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}