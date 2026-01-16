package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if cfg.Export.OutputDir != "" {
		t.Errorf("default Export.OutputDir should be empty, got %q", cfg.Export.OutputDir)
	}

	if len(cfg.Tags.TaskType) != 0 {
		t.Errorf("default Tags.TaskType should be empty, got %v", cfg.Tags.TaskType)
	}
}

func TestTemplateConfig(t *testing.T) {
	cfg := TemplateConfig()

	if cfg == nil {
		t.Fatal("TemplateConfig() returned nil")
	}

	if len(cfg.Tags.TaskType) == 0 {
		t.Error("TemplateConfig() Tags.TaskType should not be empty")
	}

	if cfg.Export.OutputDir == "" {
		t.Error("TemplateConfig() Export.OutputDir should not be empty")
	}
}

func TestLoadFromPath_NotExists(t *testing.T) {
	cfg, err := LoadFromPath("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("LoadFromPath() should not error for non-existent file, got %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadFromPath() returned nil for non-existent file")
	}

	// Should return default config
	if cfg.Export.OutputDir != "" {
		t.Errorf("should return default config, got OutputDir = %q", cfg.Export.OutputDir)
	}
}

func TestLoadFromPath_ValidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
tags:
  task_type:
    - 新機能
    - バグ修正
  project:
    - プロジェクトA
export:
  output_dir: /tmp/export
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("LoadFromPath() error = %v", err)
	}

	if len(cfg.Tags.TaskType) != 2 {
		t.Errorf("Tags.TaskType length = %d, want 2", len(cfg.Tags.TaskType))
	}

	if cfg.Tags.TaskType[0] != "新機能" {
		t.Errorf("Tags.TaskType[0] = %q, want %q", cfg.Tags.TaskType[0], "新機能")
	}

	if cfg.Export.OutputDir != "/tmp/export" {
		t.Errorf("Export.OutputDir = %q, want %q", cfg.Export.OutputDir, "/tmp/export")
	}
}

func TestLoadFromPath_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	invalidYAML := `
tags:
  - this is invalid
  task_type: [
`
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := LoadFromPath(configPath)
	if err == nil {
		t.Error("LoadFromPath() should error for invalid YAML")
	}
}

func TestConfig_SaveToPath(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.yaml")

	cfg := TemplateConfig()
	if err := cfg.SaveToPath(configPath); err != nil {
		t.Fatalf("SaveToPath() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}

	// Load and verify
	loaded, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("LoadFromPath() error = %v", err)
	}

	if len(loaded.Tags.TaskType) != len(cfg.Tags.TaskType) {
		t.Errorf("loaded Tags.TaskType length = %d, want %d", len(loaded.Tags.TaskType), len(cfg.Tags.TaskType))
	}
}

func TestConfig_OutputDir(t *testing.T) {
	tests := []struct {
		name      string
		outputDir string
		wantEmpty bool
	}{
		{
			name:      "empty",
			outputDir: "",
			wantEmpty: true,
		},
		{
			name:      "absolute path",
			outputDir: "/tmp/export",
			wantEmpty: false,
		},
		{
			name:      "tilde path",
			outputDir: "~/Documents",
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Export: ExportConfig{
					OutputDir: tt.outputDir,
				},
			}

			got := cfg.OutputDir()
			if tt.wantEmpty && got != "" {
				t.Errorf("OutputDir() = %q, want empty", got)
			}
			if !tt.wantEmpty && got == "" {
				t.Errorf("OutputDir() = empty, want non-empty")
			}

			// Check tilde expansion
			if tt.outputDir == "~/Documents" {
				if got[0] == '~' {
					t.Error("OutputDir() should expand ~")
				}
			}
		})
	}
}

func TestConfig_GetExportOutputDir_EnvOverride(t *testing.T) {
	cfg := &Config{
		Export: ExportConfig{
			OutputDir: "/config/path",
		},
	}

	// Without env var
	os.Unsetenv("DJOU_EXPORT_DIR")
	got := cfg.GetExportOutputDir()
	if got != "/config/path" {
		t.Errorf("GetExportOutputDir() = %q, want %q", got, "/config/path")
	}

	// With env var
	os.Setenv("DJOU_EXPORT_DIR", "/env/path")
	defer os.Unsetenv("DJOU_EXPORT_DIR")

	got = cfg.GetExportOutputDir()
	if got != "/env/path" {
		t.Errorf("GetExportOutputDir() with env = %q, want %q", got, "/env/path")
	}
}

func TestTemplateYAML(t *testing.T) {
	yaml := TemplateYAML()

	if yaml == "" {
		t.Error("TemplateYAML() should not return empty string")
	}

	// Check for expected content
	expectedContents := []string{
		"tags:",
		"task_type:",
		"export:",
		"output_dir:",
		"DJOU_EXPORT_DIR",
	}

	for _, expected := range expectedContents {
		if !contains(yaml, expected) {
			t.Errorf("TemplateYAML() should contain %q", expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
