package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempDockerConfig(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return configPath
}

func TestGetCredentialsFromDockerConfig(t *testing.T) {
	username := "testuser"
	password := "testpass"
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	badAuth := base64.StdEncoding.EncodeToString([]byte("no-colon"))

	tests := []struct {
		name     string
		config   string
		wantUser string
		wantPass string
		wantErr  string
	}{
		{
			name: "When a valid config is passed then the correct credentials are returned",
			config: `{
				"auths": {
					"https://index.docker.io/v1/": {
						"auth": "` + auth + `"
					}
				}
			}`,
			wantUser: username,
			wantPass: password,
			wantErr:  "",
		},
		{
			name:    "When no auths are present then an error is returned",
			config:  `{"auths":{}}`,
			wantErr: "no auth config found for registry: https://index.docker.io/v1/",
		},
		{
			name: "When invalid base64 is present then an error is returned",
			config: `{
				"auths": {
					"https://index.docker.io/v1/": {
						"auth": "not-base64"
					}
				}
			}`,
			wantErr: "illegal base64 data", // partial match
		},
		{
			name: "When invalid format is present then an error is returned",
			config: `{
				"auths": {
					"https://index.docker.io/v1/": {
						"auth": "` + badAuth + `"
					}
				}
			}`,
			wantErr: "invalid auth format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := writeTempDockerConfig(t, tt.config)
			u, p, err := getCredentialsFromDockerConfig(configPath)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if u != tt.wantUser || p != tt.wantPass {
					t.Errorf("expected %s/%s, got %s/%s", tt.wantUser, tt.wantPass, u, p)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected error containing %q, got %v", tt.wantErr, err)
				}
			}
		})
	}

	t.Run("file not found", func(t *testing.T) {
		_, _, err := getCredentialsFromDockerConfig("/nonexistent/path/config.json")
		if err == nil {
			t.Errorf("expected file not found error, got nil")
		}
	})
}
