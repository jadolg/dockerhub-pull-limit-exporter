package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getToken(username, password string, timeout time.Duration) (string, error) {
	url := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch token: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	token, ok := result["token"].(string)
	if !ok {
		return "", errors.New("token not found in response")
	}

	return token, nil
}

func getLimits(token string, timeout time.Duration) (int, int, int, int, string, error) {
	url := "https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest"

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, 0, 0, 0, "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, 0, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, 0, 0, "", fmt.Errorf("failed to fetch limits: status code %d", resp.StatusCode)
	}

	limit := resp.Header.Get("ratelimit-limit")
	remaining := resp.Header.Get("ratelimit-remaining")
	source := resp.Header.Get("docker-ratelimit-source")

	limitInt, limitWindow, remainingInt, remainingWindow, err := parseLimits(limit, remaining)
	if err != nil {
		return 0, 0, 0, 0, "", err
	}

	return limitInt, remainingInt, limitWindow, remainingWindow, source, nil
}

func parseLimits(limit string, remaining string) (int, int, int, int, error) {
	limitParts := strings.Split(limit, ";")
	limitInt, err := strconv.Atoi(limitParts[0])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to parse ratelimit-limit: %w", err)
	}
	limitWindow := 0
	if len(limitParts) > 1 {
		windowParts := strings.Split(limitParts[1], "=")
		if len(windowParts) > 1 {
			limitWindow, err = strconv.Atoi(windowParts[1])
			if err != nil {
				return 0, 0, 0, 0, fmt.Errorf("failed to parse ratelimit-limit window: %w", err)
			}
		}
	} else {
		return 0, 0, 0, 0, errors.New("ratelimit-limit header does not contain window information")
	}

	remainingParts := strings.Split(remaining, ";")
	remainingInt, err := strconv.Atoi(remainingParts[0])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to parse ratelimit-remaining: %w", err)
	}
	remainingWindow := 0
	if len(remainingParts) > 1 {
		windowParts := strings.Split(remainingParts[1], "=")
		if len(windowParts) > 1 {
			remainingWindow, err = strconv.Atoi(windowParts[1])
			if err != nil {
				return 0, 0, 0, 0, fmt.Errorf("failed to parse ratelimit-remaining window: %w", err)
			}
		}
	} else {
		return 0, 0, 0, 0, errors.New("ratelimit-remaining header does not contain window information")
	}
	return limitInt, limitWindow, remainingInt, remainingWindow, nil
}
