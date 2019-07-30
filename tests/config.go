package tests

import (
	"os"
	"strconv"
)

type Config struct {
	ProxyServerPath string
	MaxRatePerIP    uint64
	TotalMaxRate    uint64
}

var config Config

func LoadConfig() Config {
	if config.ProxyServerPath != "" {
		return config
	}

	config = Config{}

	proxyServerPath := os.Getenv("PROXY_SERVER_PATH")
	if proxyServerPath == "" {
		proxyServerPath = "http://localhost:9020"
	}

	maxRatePerIP, err := strconv.ParseUint(os.Getenv("MAX_RATE_PER_IP"), 10, 64)
	if err != nil {
		maxRatePerIP = 100
	}

	totalMaxRate, err := strconv.ParseUint(os.Getenv("TOTAL_MAX_RATE"), 10, 64)
	if err != nil {
		totalMaxRate = 10000
	}

	config.ProxyServerPath = proxyServerPath
	config.MaxRatePerIP = maxRatePerIP
	config.TotalMaxRate = totalMaxRate

	return config
}
