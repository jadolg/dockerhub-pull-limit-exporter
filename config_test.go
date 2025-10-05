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

	if config.Credentials[0].Username != "user1" {
		t.Fatalf("Expected username to be 'user1', got '%s'", config.Credentials[0].Username)
	}
	if config.Credentials[0].Password != "password1" {
		t.Fatalf("Expected password to be 'password1', got '%s'", config.Credentials[0].Password)
	}

	if config.Credentials[2].Username != "user001" {
		t.Fatalf("Expected username to be 'user001', got '%s'", config.Credentials[2].Username)
	}
	if config.Credentials[2].Password != "asimplepassword" {
		t.Fatalf("Expected password to be 'asimplepassword', got '%s'", config.Credentials[2].Password)
	}
}
