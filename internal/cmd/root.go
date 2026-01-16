package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "djou",
	Short: "Dev Journal - 開発日誌CLI",
	Long:  "日々の開発作業を記録・管理するCLIツール",
	Run: func(cmd *cobra.Command, args []string) {
		// Default: run interactive log entry
		fmt.Println("djou - Dev Journal")
		fmt.Println("対話形式での記録は今後実装予定")
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
