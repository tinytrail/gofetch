package config

import (
	"testing"
)

func TestGetDefaultAddress(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "no environment variable",
			envValue: "",
			expected: DefaultPort,
		},
		{
			name:     "valid port",
			envValue: "9090",
			expected: ":9090",
		},
		{
			name:     "invalid port - non-numeric",
			envValue: "abc",
			expected: DefaultPort,
		},
		{
			name:     "invalid port - out of range",
			envValue: "70000",
			expected: DefaultPort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				t.Setenv("MCP_PORT", tt.envValue)
			}

			result := GetDefaultAddress()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
