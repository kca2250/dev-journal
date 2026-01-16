package cmd

import (
	"github.com/kca2250/djou/internal/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP サーバーとして起動",
	Long:  "Model Context Protocol (MCP) サーバーとして起動し、Claude Code などの MCP クライアントから利用可能にします。",
	RunE:  runMCP,
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}

func runMCP(cmd *cobra.Command, args []string) error {
	return mcp.Run()
}
