// Package config provides an abstraction for the project's configurable items
package config

import (
	"os"
	"strconv"
	"strings"
)

type Exchange struct {
	BaseURL       string   `json:"base_url"`
	WindowSize    int      `json:"window_size"`
	Channels      []string `json:"channels"`
	Subscriptions []string `json:"subscriptions"`
}

type Config struct {
	Exchange  Exchange
	DebugPort string `json:"debug_port"`
}

func New() *Config {
	return &Config{
		Exchange: Exchange{
			BaseURL:       getEnv("COINBASE_BASE_URL", ""),
			WindowSize:    getEnvAsInt("WINDOW_SIZE", 10),
			Channels:      getEnvAsArray("CHANNELS", ","),
			Subscriptions: getEnvAsArray("SUBSCRIPTIONS", ","),
		},
		DebugPort: getEnv("DEBUG_PORT", "12000"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsArray(name string, sep string) []string {
	values := getEnv(name, "")
	val := strings.Split(values, sep)
	return val
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
