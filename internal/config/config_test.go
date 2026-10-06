package config

import (
	"os"
	"testing"
)

// Test the Load function to ensure it correctly loads configuration from environment variables and flags.
func TestLoad(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("PORT", "9090")
	os.Setenv("WORKER_COUNT", "5")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("WORKER_COUNT")

	// Load configuration
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate loaded configuration
	if cfg.Port != "9090" {
		t.Errorf("Expected Port to be '9090', got '%s'", cfg.Port)
	}
	if cfg.WorkerCount != 5 {
		t.Errorf("Expected WorkerCount to be 5, got %d", cfg.WorkerCount)
	}
}

// Test that invalid values (e.g. WorkerCount: 0) trigger a validation error!
func TestLoad_InvalidWorkerCount(t *testing.T) {
	// Set environment variable for testing
	os.Setenv("WORKER_COUNT", "0")
	defer os.Unsetenv("WORKER_COUNT")

	// Load configuration
	_, err := Load()
	if err == nil {
		t.Fatalf("Expected error due to invalid WorkerCount, got nil")
	}
}

func TestLoad_FlagOverrides(t *testing.T) {
	cfg, err := Load("-port", "7070", "-workers", "8")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cfg.Port != "7070" {
		t.Errorf("Expected Port to be '7070', got '%s'", cfg.Port)
	}
	if cfg.WorkerCount != 8 {
		t.Errorf("Expected WorkerCount to be 8, got %d", cfg.WorkerCount)
	}
}
