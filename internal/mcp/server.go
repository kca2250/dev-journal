package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

const (
	ServerName    = "djou"
	ServerVersion = "1.0.0"
)

// NewServer creates a new MCP server with all djou tools registered
func NewServer() *server.MCPServer {
	s := server.NewMCPServer(
		ServerName,
		ServerVersion,
		server.WithToolCapabilities(false),
	)

	// Register tools
	registerRecordTool(s)
	registerListTool(s)
	registerSearchTool(s)
	registerStatsTool(s)
	registerExportTool(s)

	return s
}

// Run starts the MCP server on stdio
func Run() error {
	s := NewServer()
	return server.ServeStdio(s)
}
