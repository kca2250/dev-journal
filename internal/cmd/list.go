package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kca2250/djou/internal/db"
	"github.com/kca2250/djou/internal/model"
	"github.com/kca2250/djou/internal/ui"
	"github.com/spf13/cobra"
)

var (
	listWeek  bool
	listMonth bool
	listLimit int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "開発日誌の一覧を表示",
	Long:  "記録した開発日誌を一覧表示します。オプションで期間を絞り込めます。",
	RunE:  runList,
}

func init() {
	listCmd.Flags().BoolVarP(&listWeek, "week", "w", false, "今週の記録を表示")
	listCmd.Flags().BoolVarP(&listMonth, "month", "m", false, "今月の記録を表示")
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 10, "表示件数を指定")
	rootCmd.AddCommand(listCmd)
}

// ErrExclusiveOptions is returned when mutually exclusive options are used together
var ErrExclusiveOptions = errors.New("--week and --month cannot be used together")

func runList(cmd *cobra.Command, args []string) error {
	// Initialize localizer
	localizer := ui.NewAutoLocalizer()
	ui.SetGlobalLocalizer(localizer)

	// Check for exclusive options
	if listWeek && listMonth {
		return ErrExclusiveOptions
	}

	// Open database
	database, err := db.Open()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}
	defer database.Close()

	// Create repository
	repo := db.NewLogRepository(database)

	// Build list options
	opts := db.ListOptions{
		Limit: listLimit,
		Week:  listWeek,
	}
	if listMonth {
		opts.Month = "current"
	}

	// Fetch logs
	ctx := context.Background()
	logs, err := repo.List(ctx, opts)
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Check if no logs found
	if len(logs) == 0 {
		fmt.Println(localizer.Get(ui.MsgNoLogs))
		return nil
	}

	// Render table
	renderLogTable(os.Stdout, logs, localizer)

	return nil
}

func renderLogTable(w *os.File, logs []model.Log, l *ui.Localizer) {
	headers := []string{
		l.Get(ui.MsgHeaderDate),
		l.Get(ui.MsgHeaderTask),
		l.Get(ui.MsgHeaderEstimate),
		l.Get(ui.MsgHeaderActual),
		l.Get(ui.MsgHeaderAI),
	}

	rows := make([][]string, len(logs))
	for i, log := range logs {
		rows[i] = []string{
			log.CreatedAt.Format("2006-01-02"),
			ui.TruncateString(log.TaskName, 20),
			ui.FormatHours(log.EstimateHours),
			ui.FormatHours(log.ActualHours),
			ui.FormatAIMinutes(log.AIMinutes),
		}
	}

	ui.RenderTable(w, headers, rows)
}
