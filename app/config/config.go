package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Version string `yaml:"version"`
	AppPort string `yaml:"appPort"`
	DatabaseAddr string `yaml:"databaseAddr"`
	DatabaseName string `yaml:"databaseName"`
	DatabaseUser string `yaml:"databaseUser"`
	DatabasePassword string `yaml:"databasePassword"`
}

// GetConfig generates settings from config file
func GetConfig(path string) (Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return Config{}, err
	}

	config.AppPort = os.Getenv("PORT")

	return config, nil
}