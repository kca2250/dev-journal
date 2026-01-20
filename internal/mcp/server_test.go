package mcp

import (
	"testing"

	"github.com/kca2250/djou/internal/version"
)

func TestNewServer(t *testing.T) {
	s := NewServer()
	if s == nil {
		t.Fatal("NewServer() returned nil")
	}

	// Check that all 7 tools are registered
	tools := s.ListTools()
	expectedTools := []string{
		"djou_record",
		"djou_list",
		"djou_search",
		"djou_stats",
		"djou_export",
		"djou_update",
		"djou_delete",
	}

	for _, name := range expectedTools {
		if _, ok := tools[name]; !ok {
			t.Errorf("expected tool %q to be registered", name)
		}
	}
}

func TestServerConstants(t *testing.T) {
	if ServerName != "djou" {
		t.Errorf("ServerName = %q, want %q", ServerName, "djou")
	}
	// Version should be set via ldflags or default to "dev"
	v := version.GetVersion()
	if v == "" {
		t.Error("version.GetVersion() returned empty string")
	}
}
