package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Tags   TagConfig   `yaml:"tags,omitempty"`
	Export ExportConfig `yaml:"export,omitempty"`
}

// TagConfig holds tag definitions
type TagConfig struct {
	TaskType []string `yaml:"task_type,omitempty"`
	Project  []string `yaml:"project,omitempty"`
	TechArea []string `yaml:"tech_area,omitempty"`
}

// ExportConfig holds export settings
type ExportConfig struct {
	OutputDir string `yaml:"output_dir,omitempty"`
}

const (
	// ConfigDir is the directory name for djou configuration
	ConfigDir = ".djou"
	// ConfigFile is the config file name
	ConfigFile = "config.yaml"
)

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Tags: TagConfig{
			TaskType: []string{},
			Project:  []string{},
			TechArea: []string{},
		},
		Export: ExportConfig{
			OutputDir: "",
		},
	}
}

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ConfigDir, ConfigFile), nil
}

// ConfigDirPath returns the path to the config directory
func ConfigDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ConfigDir), nil
}

// Load loads the configuration from the config file
// If the file does not exist, it returns the default configuration
func Load() (*Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	return LoadFromPath(configPath)
}

// LoadFromPath loads the configuration from a specific path
func LoadFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// Save saves the configuration to the config file
func (c *Config) Save() error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	return c.SaveToPath(configPath)
}

// SaveToPath saves the configuration to a specific path
func (c *Config) SaveToPath(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Exists checks if the config file exists
func Exists() (bool, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetExportOutputDir returns the export output directory
// It respects environment variable DJOU_EXPORT_DIR if set
func (c *Config) GetExportOutputDir() string {
	// Environment variable takes precedence
	if envDir := os.Getenv("DJOU_EXPORT_DIR"); envDir != "" {
		return envDir
	}
	return c.OutputDir()
}

// OutputDir returns the configured output directory, expanding ~ if present
func (c *Config) OutputDir() string {
	dir := c.Export.OutputDir
	if dir == "" {
		return ""
	}

	// Expand ~ to home directory
	if len(dir) > 0 && dir[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			dir = filepath.Join(homeDir, dir[1:])
		}
	}

	return dir
}

// TemplateConfig returns a sample configuration with example values
func TemplateConfig() *Config {
	return &Config{
		Tags: TagConfig{
			TaskType: []string{"新機能", "バグ修正", "リファクタ", "調査"},
			Project:  []string{"案件A", "案件B"},
			TechArea: []string{"フロント", "バック", "インフラ"},
		},
		Export: ExportConfig{
			OutputDir: "~/Documents/dev-journal",
		},
	}
}

// TemplateYAML returns a YAML template with comments
func TemplateYAML() string {
	return `# ~/.djou/config.yaml
# djou設定ファイル

# タグ定義
tags:
  # タスクの種類
  task_type:
    - 新機能
    - バグ修正
    - リファクタ
    - 調査
  # プロジェクト名
  project:
    - 案件A
    - 案件B
  # 技術領域
  tech_area:
    - フロント
    - バック
    - インフラ

# エクスポート設定
export:
  # デフォルトの出力先ディレクトリ
  # 環境変数 DJOU_EXPORT_DIR でオーバーライド可能
  output_dir: ~/Documents/dev-journal
`
}
