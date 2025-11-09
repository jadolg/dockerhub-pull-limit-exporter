package main

import (
	"errors"
	"os"
	"testing"
	"time"
)

func GetCredentialsFromEnv() (string, string) {
	username := os.Getenv("DOCKERHUB_USERNAME")
	password := os.Getenv("DOCKERHUB_PASSWORD")
	if username == "" || password == "" {
		panic("DOCKERHUB_USERNAME and DOCKERHUB_PASSWORD must be set")
	}
	return username, password
}

func TestGetToken(t *testing.T) {
	username, password := GetCredentialsFromEnv()
	token, err := getToken(username, password, 10*time.Second)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("Expected token to be non-empty, got empty string")
	}
}

func TestGetLimits(t *testing.T) {
	username, password := GetCredentialsFromEnv()

	token, err := getToken(username, password, 10*time.Second)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	limit, remaining, limitWindow, remainingWindow, source, err := getLimits(token, 10*time.Second)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if limit == 0 || remaining == 0 || limitWindow == 0 || remainingWindow == 0 {
		t.Fatalf("Expected limit and remaining to be non-zero, got %d and %d", limit, remaining)
	}
	if source == "" {
		t.Fatalf("Expected source to be non-empty, got empty string")
	}
	t.Logf("Limit: %d, Remaining: %d, Source: %s", limit, remaining, source)
}

func TestGetLimitsNoUser(t *testing.T) {
	token, err := getToken("", "", 10*time.Second)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	limit, remaining, limitWindow, remainingWindow, source, err := getLimits(token, 10*time.Second)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if limit == 0 || remaining == 0 || limitWindow == 0 || remainingWindow == 0 {
		t.Fatalf("Expected limit and remaining to be non-zero, got %d and %d", limit, remaining)
	}
	if source == "" {
		t.Fatalf("Expected source to be non-empty, got empty string")
	}
	t.Logf("Limit: %d, Remaining: %d, Source: %s", limit, remaining, source)
}

func TestParseLimits(t *testing.T) {
	tests := []struct {
		name                    string
		limit                   string
		remaining               string
		expectedLimit           int
		expectedWindowLimit     int
		expectedRemaining       int
		expectedWindowRemaining int
		expectedError           error
	}{
		{
			name:                    "Valid headers",
			limit:                   "100;w=60",
			remaining:               "50;w=60",
			expectedLimit:           100,
			expectedWindowLimit:     60,
			expectedRemaining:       50,
			expectedWindowRemaining: 60,
			expectedError:           nil,
		},
		{
			name:                    "Missing window in limit",
			limit:                   "100",
			remaining:               "50;w=60",
			expectedLimit:           0,
			expectedWindowLimit:     0,
			expectedRemaining:       0,
			expectedWindowRemaining: 0,
			expectedError:           errors.New("ratelimit-limit header does not contain window information"),
		},
		{
			name:                    "Missing window in remaining",
			limit:                   "100;w=60",
			remaining:               "50",
			expectedLimit:           0,
			expectedWindowLimit:     0,
			expectedRemaining:       0,
			expectedWindowRemaining: 0,
			expectedError:           errors.New("ratelimit-remaining header does not contain window information"),
		},
		{
			name:                    "Invalid limit value",
			limit:                   "abc;w=60",
			remaining:               "50;w=60",
			expectedLimit:           0,
			expectedWindowLimit:     0,
			expectedRemaining:       0,
			expectedWindowRemaining: 0,
			expectedError:           errors.New("failed to parse ratelimit-limit: strconv.Atoi: parsing \"abc\": invalid syntax"),
		},
		{
			name:                    "Invalid remaining value",
			limit:                   "100;w=60",
			remaining:               "xyz;w=60",
			expectedLimit:           0,
			expectedWindowLimit:     0,
			expectedRemaining:       0,
			expectedWindowRemaining: 0,
			expectedError:           errors.New("failed to parse ratelimit-remaining: strconv.Atoi: parsing \"xyz\": invalid syntax"),
		},
		{
			name:                    "Invalid window in limit",
			limit:                   "100;w=abc",
			remaining:               "50;w=60",
			expectedLimit:           0,
			expectedWindowLimit:     0,
			expectedRemaining:       0,
			expectedWindowRemaining: 0,
			expectedError:           errors.New("failed to parse ratelimit-limit window: strconv.Atoi: parsing \"abc\": invalid syntax"),
		},
		{
			name:                    "Invalid window in remaining",
			limit:                   "100;w=60",
			remaining:               "50;w=abc",
			expectedLimit:           0,
			expectedWindowLimit:     0,
			expectedRemaining:       0,
			expectedWindowRemaining: 0,
			expectedError:           errors.New("failed to parse ratelimit-remaining window: strconv.Atoi: parsing \"abc\": invalid syntax"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, limitWindow, remaining, remainingWindow, err := parseLimits(tt.limit, tt.remaining)

			if limit != tt.expectedLimit {
				t.Errorf("expected limit %d, got %d", tt.expectedLimit, limit)
			}
			if limitWindow != tt.expectedWindowLimit {
				t.Errorf("expected limit window %d, got %d", tt.expectedWindowLimit, limitWindow)
			}
			if remaining != tt.expectedRemaining {
				t.Errorf("expected remaining %d, got %d", tt.expectedRemaining, remaining)
			}
			if remainingWindow != tt.expectedWindowRemaining {
				t.Errorf("expected remaining window %d, got %d", tt.expectedWindowRemaining, remainingWindow)
			}
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
