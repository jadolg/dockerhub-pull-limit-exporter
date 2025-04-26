package main

import "testing"

func TestGetConfig(t *testing.T) {
	config, err := getConfig("config.example.yaml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(config.Credentials) == 0 {
		t.Fatalf("Expected at least one credential, got none")
	}
	for _, cred := range config.Credentials {
		if cred.Username == "" || cred.Password == "" {
			t.Fatalf("Expected non-empty username and password, got %s and %s", cred.Username, cred.Password)
		}
	}
}
