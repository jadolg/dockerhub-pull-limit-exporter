package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetConfig(t *testing.T) {
	config, err := getConfig("config.example.yaml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(config.Credentials) == 0 {
		t.Fatalf("Expected at least one credential, got none")
	}
	for _, cred := range config.Credentials {
		if (cred.Username == "" || cred.Password == "") && !cred.Anonymous {
			t.Fatalf("Expected non-empty username and password, got %s and %s", cred.Username, cred.Password)
		}
		if cred.Anonymous && (cred.Username != "" || cred.Password != "") {
			t.Fatalf("Expected empty credentials for anonymous user but got username: %s and password %s instead", cred.Username, cred.Password)
		}
	}

	tests := []struct {
		username string
		password string
	}{
		{"user1", "password1"},
		{"user2", "password2"},
		{"user001", "asimplepassword"},
		{"user002", "anothersimplepassword"},
	}
	for i, tt := range tests {
		cred := config.Credentials[i]
		if cred.Username != tt.username {
			t.Fatalf("Expected username to be '%s', got '%s'", tt.username, cred.Username)
		}
		if cred.Password != tt.password {
			t.Fatalf("Expected password to be '%s', got '%s'", tt.password, cred.Password)
		}
	}
}

func TestInvalidCredentials(t *testing.T) {
	cfg := configuration{
		Credentials: []credentials{
			{Username: "", Password: "", Anonymous: false},
		},
		UpdateInterval: time.Second,
		Timeout:        time.Second,
	}
	_, err := func() (configuration, error) {
		for _, credential := range cfg.Credentials {
			if credential.invalid() {
				return configuration{}, fmt.Errorf("invalid credentials configuration detected for user [%s]", credential.Username)
			}
		}
		return cfg, nil
	}()
	if err == nil {
		t.Fatal("expected error for invalid credentials, got nil")
	}
}

func TestMissingUpdateInterval(t *testing.T) {
	cfg := configuration{
		Credentials:    []credentials{{Username: "user", Password: "pass"}},
		UpdateInterval: 0,
		Timeout:        time.Second,
	}
	_, err := func() (configuration, error) {
		if cfg.UpdateInterval == 0 {
			return configuration{}, fmt.Errorf("update interval must be set")
		}
		return cfg, nil
	}()
	if err == nil || err.Error() != "update interval must be set" {
		t.Fatalf("expected error for missing update interval, got %v", err)
	}
}

func TestMissingTimeout(t *testing.T) {
	cfg := configuration{
		Credentials:    []credentials{{Username: "user", Password: "pass"}},
		UpdateInterval: time.Second,
		Timeout:        0,
	}
	_, err := func() (configuration, error) {
		if cfg.Timeout == 0 {
			return configuration{}, fmt.Errorf("timeout must be set")
		}
		return cfg, nil
	}()
	if err == nil || err.Error() != "timeout must be set" {
		t.Fatalf("expected error for missing timeout, got %v", err)
	}
}
