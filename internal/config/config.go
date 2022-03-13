package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Exchange struct {
	BaseUrl string `json:"base_url"`
	WindowSize int `json:"window_size"`
	Channels []string `json:"channels"`
	Subscriptions []string `json:"subscriptions"`
}

type Config struct {
	Exchange Exchange
	DebugPort string `json:"debug_port"`
}

func Parse() (*Config, error) {
	return fileToConfig("config/config.json")
}

func FromPath(url string) (*Config, error) {
	return fileToConfig(url)
}

func fileToConfig(path string) (*Config, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	rawFile, _ := ioutil.ReadAll(jsonFile)
	var config Config

	if err = json.Unmarshal(rawFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}