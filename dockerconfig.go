package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func getCredentialsFromDockerConfig(configPath string) (string, string, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return "", "", err
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Errorf("Error closing config file: %v", err)
		}
	}(configFile)
	var config struct {
		Auths map[string]struct {
			Auth string `json:"auth"`
		} `json:"auths"`
	}
	bytes, err := io.ReadAll(configFile)
	if err != nil {
		return "", "", err
	}
	if err := json.Unmarshal(bytes, &config); err != nil {
		return "", "", err
	}
	registryURL := "https://index.docker.io/v1/"

	authEntry, ok := config.Auths[registryURL]
	if !ok {
		return "", "", fmt.Errorf("no auth config found for registry: %s", registryURL)
	}
	decoded, err := base64.StdEncoding.DecodeString(authEntry.Auth)
	if err != nil {
		return "", "", err
	}
	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return "", "", fmt.Errorf("invalid auth format")
	}

	return strings.TrimSpace(credentials[0]), strings.TrimSpace(credentials[1]), nil
}
