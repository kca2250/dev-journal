package cmd

import (
	"fmt"

	"github.com/kca2250/djou/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョン情報を表示",
	Long:  `djouのバージョン情報を表示します。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("djou " + version.GetFullVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
