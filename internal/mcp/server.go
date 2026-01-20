package mcp

import (
	"github.com/kca2250/djou/internal/version"
	"github.com/mark3labs/mcp-go/server"
)

const (
	ServerName = "djou"
)

// NewServer creates a new MCP server with all djou tools registered
func NewServer() *server.MCPServer {
	s := server.NewMCPServer(
		ServerName,
		version.GetVersion(),
		server.WithToolCapabilities(false),
	)

	// Register tools
	registerRecordTool(s)
	registerListTool(s)
	registerSearchTool(s)
	registerStatsTool(s)
	registerExportTool(s)
	registerUpdateTool(s)
	registerDeleteTool(s)

	return s
}

// Run starts the MCP server on stdio
func Run() error {
	s := NewServer()
	return server.ServeStdio(s)
}
