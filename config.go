package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type configuration struct {
	Credentials    []credentials `json:"credentials"`
	UpdateInterval time.Duration `yaml:"update_interval"`
	Timeout        time.Duration `yaml:"timeout"`
	ConfigFiles    []string      `yaml:"config_files"`
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c credentials) invalid() bool {
	return c.Username == "" || c.Password == ""
}

func getConfig(configFile string) (configuration, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return configuration{}, err
	}
	c := configuration{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return configuration{}, err
	}

	for _, configFile := range c.ConfigFiles {
		username, password, err := getCredentialsFromDockerConfig(configFile)
		if err != nil {
			return configuration{}, fmt.Errorf("error reading docker config file %s: %v", configFile, err)
		}
		if username == "" || password == "" {
			return configuration{}, fmt.Errorf("invalid docker config file detected %s", configFile)
		}

		c.Credentials = append(c.Credentials, credentials{
			Username: username,
			Password: password,
		})
	}

	for _, credential := range c.Credentials {
		if credential.invalid() {
			return configuration{}, fmt.Errorf("invalid credentials configuration detected %+v", credential)
		}
	}

	if c.UpdateInterval == 0 {
		return configuration{}, fmt.Errorf("update interval must be set")
	}

	if c.Timeout == 0 {
		return configuration{}, fmt.Errorf("timeout must be set")
	}

	return c, nil
}
