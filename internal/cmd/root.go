package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kca2250/djou/internal/db"
	"github.com/kca2250/djou/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "djou",
	Short: "Dev Journal - 開発日誌CLI",
	Long:  "日々の開発作業を記録・管理するCLIツール",
	RunE:  runRecord,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runRecord(cmd *cobra.Command, args []string) error {
	// Initialize localizer
	localizer := ui.NewAutoLocalizer()
	ui.SetGlobalLocalizer(localizer)

	// Open database
	database, err := db.Open()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}
	defer database.Close()

	// Create repository
	repo := db.NewLogRepository(database)

	// Run the form
	form := ui.NewRecordForm(localizer)
	input, err := form.Run()
	if err != nil {
		if errors.Is(err, ui.ErrFormCancelled) {
			fmt.Println(localizer.Get(ui.MsgRecordCancelled))
			return nil
		}
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Save to database
	ctx := context.Background()
	if err := repo.Create(ctx, input); err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Show success message
	fmt.Println("✅ " + localizer.Get(ui.MsgRecordSuccess))

	return nil
}
