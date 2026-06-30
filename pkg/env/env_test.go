package env

import (
	"os"
	"reshell/pkg/config"
	"testing"
)

func TestEnvSaveLoadToggleRemove(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "reshell-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	err = AddOrUpdate("TEST_ENV_VAR", "value123", "test variable", true)
	if err != nil {
		t.Fatalf("AddOrUpdate failed: %v", err)
	}

	cfg, err := config.LoadEnv()
	if err != nil {
		t.Fatalf("LoadEnv failed: %v", err)
	}

	if len(cfg.Variables) != 1 {
		t.Fatalf("Expected 1 env var, got %d", len(cfg.Variables))
	}

	v := cfg.Variables[0]
	if v.Name != "TEST_ENV_VAR" || v.Value != "value123" || !v.Enabled {
		t.Errorf("Variable mismatch: %+v", v)
	}

	err = Toggle("TEST_ENV_VAR")
	if err != nil {
		t.Fatalf("Toggle failed: %v", err)
	}

	cfg, _ = config.LoadEnv()
	if cfg.Variables[0].Enabled {
		t.Errorf("Expected variable to be disabled after toggle")
	}

	err = Remove("TEST_ENV_VAR")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	cfg, _ = config.LoadEnv()
	if len(cfg.Variables) != 0 {
		t.Errorf("Expected variable to be removed, got %d variables", len(cfg.Variables))
	}
}
