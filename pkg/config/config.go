// Package config provides server configuration functionality.
package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Constants
const (
	DefaultPort   = ":8080"
	ServerName    = "fetch-server"
	ServerVersion = "1.0.0"
	DefaultUA     = "Mozilla/5.0 (compatible; MCPFetchBot/1.0)"
)

// Config holds the server configuration
type Config struct {
	Address      string
	UserAgent    string
	IgnoreRobots bool
	ProxyURL     string
}

// ParseFlags parses command line flags and returns configuration
func ParseFlags() Config {
	var config Config

	addr := flag.String("addr", GetDefaultAddress(), "Address to listen on")
	flag.StringVar(&config.UserAgent, "user-agent", "", "Custom User-Agent string")
	flag.BoolVar(&config.IgnoreRobots, "ignore-robots-txt", false, "Ignore robots.txt rules")
	flag.StringVar(&config.ProxyURL, "proxy-url", "", "Proxy URL for requests")
	flag.Parse()

	config.Address = *addr

	// Set default user agent if not provided
	if config.UserAgent == "" {
		config.UserAgent = DefaultUA
	}

	return config
}

// GetDefaultAddress returns the default server address from environment or constant
func GetDefaultAddress() string {
	portEnv := os.Getenv("MCP_PORT")
	if portEnv == "" {
		return DefaultPort
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Printf("Invalid port number in MCP_PORT environment variable: %v, using default port 8080", err)
		return DefaultPort
	}

	// Validate port range
	if port < 1 || port > 65535 {
		log.Printf("Port %d out of valid range (1-65535), using default port", port)
		return DefaultPort
	}

	return fmt.Sprintf(":%d", port)
}
