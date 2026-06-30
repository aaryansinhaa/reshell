package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDirectories(t *testing.T) {
	// Set up temporary home directory for testing
	tempHome, err := os.MkdirTemp("", "reshell-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	// Ensure directories can be created
	err = EnsureDirectories()
	if err != nil {
		t.Fatalf("EnsureDirectories failed: %v", err)
	}

	configDir := filepath.Join(tempHome, ".config", "reshell")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Expected config dir %s to exist", configDir)
	}

	funcsDir := filepath.Join(configDir, "functions")
	if _, err := os.Stat(funcsDir); os.IsNotExist(err) {
		t.Errorf("Expected functions dir %s to exist", funcsDir)
	}
}

func TestAliasConfigSaveLoad(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "reshell-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	aliases := &AliasConfig{
		Aliases: []Alias{
			{Name: "gs", Value: "git status", Description: "git status alias", Shell: "all", Enabled: true},
			{Name: "la", Value: "ls -la", Description: "ls list all", Shell: "bash", Enabled: false},
		},
	}

	err = SaveAliases(aliases)
	if err != nil {
		t.Fatalf("SaveAliases failed: %v", err)
	}

	loaded, err := LoadAliases()
	if err != nil {
		t.Fatalf("LoadAliases failed: %v", err)
	}

	if len(loaded.Aliases) != 2 {
		t.Fatalf("Expected 2 aliases, got %d", len(loaded.Aliases))
	}

	if loaded.Aliases[0].Name != "gs" || loaded.Aliases[0].Value != "git status" {
		t.Errorf("First alias mismatch: %+v", loaded.Aliases[0])
	}

	if loaded.Aliases[1].Enabled {
		t.Errorf("Second alias should be disabled")
	}
}
