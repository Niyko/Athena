package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

type Config struct {
	DBHost string `json:"dbHost"`
	DBPort int    `json:"dbPort"`
	DBUser string `json:"dbUser"`
	DBPassword string `json:"dbPassword"`
	DBName string `json:"dbName"`

	KafkaHost string `json:"kafkaHost"`
	KafkaTopic string `json:"kafkaTopic"`
	KafkaSASLMechanisms string `json:"kafkaSASLMechanisms"`
	KafkaSecurityProtocol string `json:"kafkaSecurityProtocol"`
	KafkaSASLUsername string `json:"kafkaSASLUsername"`
	KafkaSASLPassword string `json:"kafkaSASLPassword"`

	PollInterval int `json:"pollInterval"`
	FetchLimit int `json:"fetchLimit"`
	SkippedTables []string `json:"skippedTables"`
}

func getConfig() Config {
	configFilePath := getExePath() + "config.json"

	configFileContent, error := ioutil.ReadFile(configFilePath)
	if error != nil {
		color.Red("Error while reading config file (%s)", error)
		os.Exit(0)
	}

	var config Config

	error = json.Unmarshal(configFileContent, &config)
	if error != nil {
		color.Red("Error while parsing config file (%s)", error)
		os.Exit(0)
	}

	return config
}

func contains[T comparable](arr []T, x T) bool {
    for _, v := range arr {
        if v == x {
            return true
        }
    }
    return false
}

func getExePath() string {
	if os.Getenv("GORUN") != "" {
		exeDir, error := os.Getwd()
		if error != nil {
			color.Red("Error while getting executable path (%s)", error)
			os.Exit(0)
		}
		
		return exeDir + string(filepath.Separator)
	} else {
		exePath, error := os.Executable()
		if error != nil {
			color.Red("Error while getting executable path (%s)", error)
			os.Exit(0)
		}
	
		exeDir := filepath.Dir(exePath)

		return exeDir + string(filepath.Separator)
	}
}