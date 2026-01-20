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
	listWeek          bool
	listMonth         bool
	listLimit         int
	listNoInteractive bool
	listTag           string
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
	listCmd.Flags().BoolVar(&listNoInteractive, "no-interactive", false, "インタラクティブモードを無効化")
	listCmd.Flags().StringVar(&listTag, "tag", "", "タグでフィルタリング")
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
		Tag:   listTag,
	}
	if listMonth {
		opts.Month = "current"
	}

	ctx := context.Background()

	// Non-interactive mode: just show the table
	if listNoInteractive {
		logs, err := repo.List(ctx, opts)
		if err != nil {
			return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
		}
		if len(logs) == 0 {
			fmt.Println(localizer.Get(ui.MsgNoLogs))
			return nil
		}
		renderLogTable(os.Stdout, logs, localizer)
		return nil
	}

	// Interactive mode: loop until user exits
	return runInteractiveList(ctx, repo, opts, localizer)
}

func runInteractiveList(ctx context.Context, repo *db.LogRepository, opts db.ListOptions, localizer *ui.Localizer) error {
	for {
		// Fetch logs
		logs, err := repo.List(ctx, opts)
		if err != nil {
			return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
		}

		if len(logs) == 0 {
			fmt.Println(localizer.Get(ui.MsgNoLogs))
			return nil
		}

		// Show log selector
		selector := ui.NewLogSelector(localizer, logs)
		selectedLog, err := selector.Run()
		if err != nil {
			return err
		}

		// User selected exit
		if selectedLog == nil {
			return nil
		}

		// Show action selector
		actionSelector := ui.NewActionSelector(localizer, *selectedLog)
		action, err := actionSelector.Run()
		if err != nil {
			return err
		}

		switch action {
		case ui.ActionEdit:
			if err := handleEdit(ctx, repo, *selectedLog, localizer); err != nil {
				if errors.Is(err, ui.ErrFormCancelled) {
					fmt.Println(localizer.Get(ui.MsgEditCancelled))
					continue
				}
				return err
			}
			fmt.Println(localizer.Get(ui.MsgEditSuccess))

		case ui.ActionDelete:
			if err := handleDelete(ctx, repo, *selectedLog, localizer); err != nil {
				return err
			}

		case ui.ActionCancel:
			// Continue to next iteration
		}
	}
}

func handleEdit(ctx context.Context, repo *db.LogRepository, log model.Log, localizer *ui.Localizer) error {
	editForm := ui.NewEditForm(localizer, log)
	input, err := editForm.Run()
	if err != nil {
		return err
	}

	return repo.Update(ctx, log.ID, input)
}

func handleDelete(ctx context.Context, repo *db.LogRepository, log model.Log, localizer *ui.Localizer) error {
	confirm := ui.NewDeleteConfirm(localizer, log)
	confirmed, err := confirm.Run()
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println(localizer.Get(ui.MsgDeleteCancelled))
		return nil
	}

	if err := repo.Delete(ctx, log.ID); err != nil {
		return err
	}

	fmt.Println(localizer.Get(ui.MsgDeleteSuccess))
	return nil
}

func renderLogTable(w *os.File, logs []model.Log, l *ui.Localizer) {
	headers := []string{
		l.Get(ui.MsgHeaderDate),
		l.Get(ui.MsgHeaderTask),
		l.Get(ui.MsgHeaderEstimate),
		l.Get(ui.MsgHeaderActual),
	}

	rows := make([][]string, len(logs))
	for i, log := range logs {
		rows[i] = []string{
			log.CreatedAt.Format("2006-01-02"),
			ui.TruncateString(log.TaskName, 20),
			ui.FormatHours(log.EstimateHours),
			ui.FormatHours(log.ActualHours),
		}
	}

	ui.RenderTable(w, headers, rows)
}
