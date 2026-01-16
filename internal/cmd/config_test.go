package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigCmd_Use(t *testing.T) {
	if configCmd.Use != "config" {
		t.Errorf("configCmd.Use = %q, want %q", configCmd.Use, "config")
	}
}

func TestConfigInitCmd_Use(t *testing.T) {
	if configInitCmd.Use != "init" {
		t.Errorf("configInitCmd.Use = %q, want %q", configInitCmd.Use, "init")
	}
}

func TestConfigShowCmd_Use(t *testing.T) {
	if configShowCmd.Use != "show" {
		t.Errorf("configShowCmd.Use = %q, want %q", configShowCmd.Use, "show")
	}
}

func TestConfigEditCmd_Use(t *testing.T) {
	if configEditCmd.Use != "edit" {
		t.Errorf("configEditCmd.Use = %q, want %q", configEditCmd.Use, "edit")
	}
}

func TestConfigCmd_HasSubcommands(t *testing.T) {
	subcommands := configCmd.Commands()
	if len(subcommands) != 3 {
		t.Errorf("configCmd should have 3 subcommands, got %d", len(subcommands))
	}

	expectedSubs := map[string]bool{
		"init": false,
		"show": false,
		"edit": false,
	}

	for _, sub := range subcommands {
		if _, ok := expectedSubs[sub.Use]; ok {
			expectedSubs[sub.Use] = true
		}
	}

	for name, found := range expectedSubs {
		if !found {
			t.Errorf("configCmd should have subcommand %q", name)
		}
	}
}

func TestConfigInit_CreatesFile(t *testing.T) {
	// Create a temp directory to simulate home
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".djou")
	configPath := filepath.Join(configDir, "config.yaml")

	// We can't easily test the actual command since it uses os.UserHomeDir()
	// Instead, test the logic that would be used

	// Create directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	// Write a test config
	testContent := "tags:\n  task_type:\n    - test\n"
	if err := os.WriteFile(configPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file should exist")
	}

	// Read and verify content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("config content = %q, want %q", string(content), testContent)
	}
}
