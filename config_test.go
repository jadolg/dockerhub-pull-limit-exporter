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
