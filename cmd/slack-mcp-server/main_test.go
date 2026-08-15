package main

import (
	"os"
	"testing"
)

func TestGetHost(t *testing.T) {
	// Backup original env vars
	origHost := os.Getenv("SLACK_MCP_HOST")
	origRailwayEnv := os.Getenv("RAILWAY_ENVIRONMENT")
	origRailwayStatic := os.Getenv("RAILWAY_STATIC_URL")
	origPort := os.Getenv("PORT")
	defer func() {
		os.Setenv("SLACK_MCP_HOST", origHost)
		os.Setenv("RAILWAY_ENVIRONMENT", origRailwayEnv)
		os.Setenv("RAILWAY_STATIC_URL", origRailwayStatic)
		os.Setenv("PORT", origPort)
	}()

	os.Unsetenv("SLACK_MCP_HOST")
	os.Unsetenv("RAILWAY_ENVIRONMENT")
	os.Unsetenv("RAILWAY_STATIC_URL")
	os.Unsetenv("PORT")

	// Default when nothing set
	if got := getHost(); got != defaultSseHost {
		t.Errorf("getHost() = %v, want %v", got, defaultSseHost)
	}

	// Explicit SLACK_MCP_HOST
	os.Setenv("SLACK_MCP_HOST", "192.168.1.1")
	if got := getHost(); got != "192.168.1.1" {
		t.Errorf("getHost() = %v, want %v", got, "192.168.1.1")
	}

	// Railway env
	os.Unsetenv("SLACK_MCP_HOST")
	os.Setenv("PORT", "8080")
	if got := getHost(); got != "0.0.0.0" {
		t.Errorf("getHost() with PORT set = %v, want %v", got, "0.0.0.0")
	}
}

func TestGetPort(t *testing.T) {
	origSlackPort := os.Getenv("SLACK_MCP_PORT")
	origPort := os.Getenv("PORT")
	defer func() {
		os.Setenv("SLACK_MCP_PORT", origSlackPort)
		os.Setenv("PORT", origPort)
	}()

	os.Unsetenv("SLACK_MCP_PORT")
	os.Unsetenv("PORT")

	// Default port
	if got := getPort(); got != "13080" {
		t.Errorf("getPort() = %v, want %v", got, "13080")
	}

	// PORT fallback
	os.Setenv("PORT", "5000")
	if got := getPort(); got != "5000" {
		t.Errorf("getPort() with PORT set = %v, want %v", got, "5000")
	}

	// SLACK_MCP_PORT precedence
	os.Setenv("SLACK_MCP_PORT", "9000")
	if got := getPort(); got != "9000" {
		t.Errorf("getPort() with SLACK_MCP_PORT set = %v, want %v", got, "9000")
	}
}
