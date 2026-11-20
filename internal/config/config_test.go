package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("TOKEN", "test-token")
	os.Setenv("CLOUD_ID", "test-cloud-id")
	defer os.Unsetenv("TOKEN")
	defer os.Unsetenv("CLOUD_ID")

	// Set command-line flags for the test
	os.Args = []string{"cmd", "-port=9090", "-host=127.0.0.1", "-timeout=10", "-debug=true"}

	cfg, err := New()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Port != "9090" {
		t.Errorf("Expected port to be '9090', got %s", cfg.Port)
	}

	if cfg.Host != "127.0.0.1" {
		t.Errorf("Expected host to be '127.0.0.1', got %s", cfg.Host)
	}

	if cfg.Timeout != 10 {
		t.Errorf("Expected timeout to be 10, got %d", cfg.Timeout)
	}

	if cfg.Debug != true {
		t.Errorf("Expected debug to be true, got %v", cfg.Debug)
	}

	if cfg.Token != "test-token" {
		t.Errorf("Expected token to be 'test-token', got %s", cfg.Token)
	}

	if cfg.CloudID != "test-cloud-id" {
		t.Errorf("Expected cloud ID to be 'test-cloud-id', got %s", cfg.CloudID)
	}
}

func TestNewMissingEnvVars(t *testing.T) {
	// Unset environment variables for the test
	os.Unsetenv("TOKEN")
	os.Unsetenv("CLOUD_ID")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to missing environment variables, but did not get one")
		}
	}()

	_, _ = New()
}
