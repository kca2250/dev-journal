package mcp

import (
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer()
	if s == nil {
		t.Fatal("NewServer() returned nil")
	}

	// Check that all 5 tools are registered
	tools := s.ListTools()
	expectedTools := []string{
		"djou_record",
		"djou_list",
		"djou_search",
		"djou_stats",
		"djou_export",
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
	if ServerVersion != "1.0.0" {
		t.Errorf("ServerVersion = %q, want %q", ServerVersion, "1.0.0")
	}
}
