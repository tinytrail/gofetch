package robots

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// createMockRobotsServer creates a test HTTP server for robots.txt testing
func createMockRobotsServer() *httptest.Server {
	mux := http.NewServeMux()

	// robots.txt endpoint
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		robotsContent := `User-agent: *
Disallow: /private/
Disallow: /admin/

User-agent: TestBot
Disallow: /blocked/`
		w.Write([]byte(robotsContent))
	})

	return httptest.NewServer(mux)
}

func TestNewChecker(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	checker := NewChecker("TestBot/1.0", false, client)

	if checker.userAgent != "TestBot/1.0" {
		t.Errorf("expected userAgent %q, got %q", "TestBot/1.0", checker.userAgent)
	}

	if checker.ignoreRobots != false {
		t.Errorf("expected ignoreRobots %v, got %v", false, checker.ignoreRobots)
	}

	if checker.httpClient != client {
		t.Error("expected httpClient to be set correctly")
	}
}

func TestIsAllowed(t *testing.T) {
	server := createMockRobotsServer()
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	tests := []struct {
		name         string
		targetURL    string
		ignoreRobots bool
		userAgent    string
		expected     bool
	}{
		{
			name:         "ignore robots enabled",
			targetURL:    server.URL + "/anything",
			ignoreRobots: true,
			userAgent:    "TestBot/1.0",
			expected:     true,
		},
		{
			name:         "allowed path",
			targetURL:    server.URL + "/public/page",
			ignoreRobots: false,
			userAgent:    "TestBot/1.0",
			expected:     true,
		},
		{
			name:         "disallowed path for all user agents",
			targetURL:    server.URL + "/private/secret",
			ignoreRobots: false,
			userAgent:    "TestBot/1.0",
			expected:     false,
		},
		{
			name:         "disallowed path for specific user agent",
			targetURL:    server.URL + "/blocked/page",
			ignoreRobots: false,
			userAgent:    "TestBot/1.0",
			expected:     false,
		},
		{
			name:         "URL with no robots.txt allows access",
			targetURL:    "http://nonexistent-host-12345.invalid/page",
			ignoreRobots: false,
			userAgent:    "TestBot/1.0",
			expected:     true, // Can't fetch robots.txt, so allow access
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewChecker(tt.userAgent, tt.ignoreRobots, client)
			result := checker.IsAllowed(tt.targetURL)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseRobotsRules(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	checker := NewChecker("TestBot/1.0", false, client)

	tests := []struct {
		name          string
		robotsContent string
		targetPath    string
		expected      bool
	}{
		{
			name: "allowed path",
			robotsContent: `User-agent: *
Disallow: /private/`,
			targetPath: "/public/page",
			expected:   true,
		},
		{
			name: "disallowed path",
			robotsContent: `User-agent: *
Disallow: /private/`,
			targetPath: "/private/secret",
			expected:   false,
		},
		{
			name: "root disallow",
			robotsContent: `User-agent: *
Disallow: /`,
			targetPath: "/anything",
			expected:   false,
		},
		{
			name: "general disallow applies to all",
			robotsContent: `User-agent: *
Disallow: /anything`,
			targetPath: "/anything",
			expected:   false,
		},
		{
			name:          "empty robots.txt",
			robotsContent: "",
			targetPath:    "/anything",
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.parseRobotsRules(tt.robotsContent, tt.targetPath)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
