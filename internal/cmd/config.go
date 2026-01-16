package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kca2250/djou/internal/config"
	"github.com/kca2250/djou/internal/ui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "設定の管理",
	Long:  "djouの設定ファイルを管理します。",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "設定ファイルを初期化",
	Long:  "~/.djou/config.yaml にテンプレート設定ファイルを作成します。",
	RunE:  runConfigInit,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "現在の設定を表示",
	Long:  "現在の設定内容を表示します。",
	RunE:  runConfigShow,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "設定ファイルを編集",
	Long:  "エディタで設定ファイルを開きます。",
	RunE:  runConfigEdit,
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configEditCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	localizer := ui.NewAutoLocalizer()
	ui.SetGlobalLocalizer(localizer)

	// Check if config already exists
	exists, err := config.Exists()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	if exists {
		fmt.Println(localizer.Get(ui.MsgConfigAlreadyExists))
		return nil
	}

	// Get config path
	configPath, err := config.ConfigPath()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Get config directory
	configDir, err := config.ConfigDirPath()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Create directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Write template YAML
	if err := os.WriteFile(configPath, []byte(config.TemplateYAML()), 0644); err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	fmt.Println(localizer.Get(ui.MsgConfigInitSuccess))
	fmt.Println(localizer.Getf(ui.MsgConfigCreated, configPath))

	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	localizer := ui.NewAutoLocalizer()
	ui.SetGlobalLocalizer(localizer)

	// Check if config exists
	exists, err := config.Exists()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	if !exists {
		fmt.Println(localizer.Get(ui.MsgConfigNotFound))
		fmt.Println(localizer.Get(ui.MsgConfigInitHint))
		return nil
	}

	// Get config path and read file
	configPath, err := config.ConfigPath()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Show config path and content
	fmt.Println(localizer.Getf(ui.MsgConfigPath, configPath))
	fmt.Println()
	fmt.Print(string(content))

	return nil
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	localizer := ui.NewAutoLocalizer()
	ui.SetGlobalLocalizer(localizer)

	// Check if config exists
	exists, err := config.Exists()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	if !exists {
		fmt.Println(localizer.Get(ui.MsgConfigNotFound))
		fmt.Println(localizer.Get(ui.MsgConfigInitHint))
		return nil
	}

	// Get editor from environment
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}

	// Get config path
	configPath, err := config.ConfigPath()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Open editor
	editorCmd := exec.Command(editor, configPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	return nil
}
