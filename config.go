package main

import (
	"encoding/json"
	"os"
)

type config struct {
	AcceptedLicenses []string          `json:"accepted-licenses"`
	Archives         map[string]string `json:"archives"`
	AlwaysInstall    []string          `json:"always-install"`
}

func readConfig(configPath string) (c config, err error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer configFile.Close()
	err = json.NewDecoder(configFile).Decode(&c)
	return
}
