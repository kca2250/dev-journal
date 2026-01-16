package cmd

import (
	"testing"
)

func TestMcpCmd_Use(t *testing.T) {
	if mcpCmd.Use != "mcp" {
		t.Errorf("mcpCmd.Use = %q, want %q", mcpCmd.Use, "mcp")
	}
}

func TestMcpCmd_Short(t *testing.T) {
	if mcpCmd.Short == "" {
		t.Error("mcpCmd.Short should not be empty")
	}
}

func TestMcpCmd_HasRunE(t *testing.T) {
	if mcpCmd.RunE == nil {
		t.Error("mcpCmd.RunE should not be nil")
	}
}
