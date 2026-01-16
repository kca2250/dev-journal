package cmd

import (
	"testing"
)

func TestRootCmd_Use(t *testing.T) {
	if rootCmd.Use != "djou" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "djou")
	}
}

func TestRootCmd_Short(t *testing.T) {
	if rootCmd.Short != "Dev Journal - 開発日誌CLI" {
		t.Errorf("rootCmd.Short = %q, want %q", rootCmd.Short, "Dev Journal - 開発日誌CLI")
	}
}

func TestRootCmd_HasRunE(t *testing.T) {
	if rootCmd.RunE == nil {
		t.Error("rootCmd.RunE should not be nil")
	}
}
