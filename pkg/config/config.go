package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppName               string
	StageStatus           string
	ServerHost            string
	ServerPort            string
	ServerDebug           bool
	ServerReadTimeout     int
	ServerShutdownTimeout int
	JourneyWorkerPool     int
	DropoffWorkerPool     int
	DropoffRetryLimit     int
}

func Load(filename string) (*Config, error) {
	var paths = []string{
		"./" + filename,
		"../" + filename,
		"../../" + filename,
		"../../../" + filename,
	}
	for _, path := range paths {
		cfg, err := New(path)
		if err == nil {
			return cfg, err
		}
	}
	return &Config{}, nil
}

func New(filename string) (*Config, error) {

	var config = &Config{}

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		var line = scanner.Text()
		var parts = strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		var key = strings.TrimSpace(parts[0])
		var value = strings.TrimSpace(parts[1])
		value = strings.ReplaceAll(value, "\"", "")

		switch key {
		case "APP_NAME":
			config.AppName = value
		case "STAGE_STATUS":
			config.StageStatus = value
		case "SERVER_HOST":
			config.ServerHost = value
		case "SERVER_PORT":
			config.ServerPort = value
		case "SERVER_DEBUG":
			if boolean, err := strconv.ParseBool(value); err == nil {
				config.ServerDebug = boolean
			}
		case "SERVER_READ_TIMEOUT":
			if number, err := strconv.Atoi(value); err == nil {
				config.ServerReadTimeout = number
			}
		case "SERVER_SHUTDOWN_TIMEOUT":
			if number, err := strconv.Atoi(value); err == nil {
				config.ServerReadTimeout = number
			}
		case "JOURNEY_WORKER_POOL":
			if number, err := strconv.Atoi(value); err == nil {
				config.JourneyWorkerPool = number
			}
		case "DROPOFF_WORKER_POOL":
			if number, err := strconv.Atoi(value); err == nil {
				config.DropoffWorkerPool = number
			}
		case "DROPOFF_RETRY_LIMIT":
			if number, err := strconv.Atoi(value); err == nil {
				config.DropoffRetryLimit = number
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return config, err
	}

	return config, nil
}
