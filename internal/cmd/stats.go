package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kca2250/djou/internal/db"
	"github.com/kca2250/djou/internal/ui"
	"github.com/spf13/cobra"
)

var statsMonth string

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "統計情報を表示",
	Long:  "記録した開発日誌の統計情報を表示します。",
	RunE:  runStats,
}

func init() {
	statsCmd.Flags().StringVarP(&statsMonth, "month", "m", "", "月別統計を表示。値を指定すると特定月のみ表示 (例: 2025-01)")
	rootCmd.AddCommand(statsCmd)
}

// ErrInvalidMonthFormat is returned when month format is invalid
var ErrInvalidMonthFormat = errors.New("invalid month format: use YYYY-MM (e.g., 2025-01)")

func runStats(cmd *cobra.Command, args []string) error {
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
	ctx := context.Background()

	// Check if --month flag was provided
	monthFlagProvided := cmd.Flags().Changed("month")

	if monthFlagProvided {
		if statsMonth == "" {
			// --month without value: show monthly stats table
			return showMonthlyStats(ctx, repo, localizer)
		}
		// --month with value: show stats for specific month
		if err := validateMonthFormat(statsMonth); err != nil {
			return err
		}
		return showStatsForMonth(ctx, repo, localizer, statsMonth)
	}

	// Default: show overall stats
	return showOverallStats(ctx, repo, localizer)
}

func showOverallStats(ctx context.Context, repo *db.LogRepository, l *ui.Localizer) error {
	stats, err := repo.GetStats(ctx)
	if err != nil {
		return fmt.Errorf("%s", l.Getf(ui.MsgError, err.Error()))
	}

	if stats.Count == 0 {
		fmt.Println(l.Get(ui.MsgStatsNoData))
		return nil
	}

	// Print stats header
	fmt.Println(l.Get(ui.MsgStatsHeader))
	fmt.Println()

	// Period
	if stats.MinDate != nil && stats.MaxDate != nil {
		fmt.Println(l.Getf(ui.MsgStatsPeriod,
			stats.MinDate.Format("2006-01-02"),
			stats.MaxDate.Format("2006-01-02")))
	}

	// Total logs
	fmt.Println(l.Getf(ui.MsgStatsTotalLogs, stats.Count))
	fmt.Println()

	// Hours
	fmt.Println(l.Getf(ui.MsgStatsTotalHours, stats.TotalEstimate, stats.TotalActual))
	fmt.Println(l.Getf(ui.MsgStatsAccuracy, calculateAccuracy(stats.TotalEstimate, stats.TotalActual)))
	fmt.Println()

	// AI usage
	aiHours := float64(stats.TotalAIMinutes) / 60.0
	aiRate := calculateAIRate(stats.TotalAIMinutes, stats.TotalActual)
	fmt.Printf(l.Getf(ui.MsgStatsTotalAI, stats.TotalAIMinutes)+" (%.1fh, %.0f%%)\n", aiHours, aiRate)

	return nil
}

func showMonthlyStats(ctx context.Context, repo *db.LogRepository, l *ui.Localizer) error {
	stats, err := repo.GetMonthlyStats(ctx)
	if err != nil {
		return fmt.Errorf("%s", l.Getf(ui.MsgError, err.Error()))
	}

	fmt.Println(l.Get(ui.MsgStatsMonthlyTable))
	fmt.Println()

	renderMonthlyStatsTable(os.Stdout, stats, l)

	return nil
}

func showStatsForMonth(ctx context.Context, repo *db.LogRepository, l *ui.Localizer, month string) error {
	stats, err := repo.GetStatsByMonth(ctx, month)
	if err != nil {
		return fmt.Errorf("%s", l.Getf(ui.MsgError, err.Error()))
	}

	if stats.Count == 0 {
		fmt.Println(l.Get(ui.MsgStatsNoData))
		return nil
	}

	// Print stats header with month
	fmt.Printf("%s (%s)\n", l.Get(ui.MsgStatsHeader), month)
	fmt.Println()

	// Total logs
	fmt.Println(l.Getf(ui.MsgStatsTotalLogs, stats.Count))
	fmt.Println()

	// Hours
	fmt.Println(l.Getf(ui.MsgStatsTotalHours, stats.TotalEstimate, stats.TotalActual))
	fmt.Println(l.Getf(ui.MsgStatsAccuracy, calculateAccuracy(stats.TotalEstimate, stats.TotalActual)))
	fmt.Println()

	// AI usage
	aiHours := float64(stats.TotalAIMinutes) / 60.0
	aiRate := calculateAIRate(stats.TotalAIMinutes, stats.TotalActual)
	fmt.Printf(l.Getf(ui.MsgStatsTotalAI, stats.TotalAIMinutes)+" (%.1fh, %.0f%%)\n", aiHours, aiRate)

	return nil
}

func renderMonthlyStatsTable(w io.Writer, stats []db.MonthlyStats, l *ui.Localizer) {
	if len(stats) == 0 {
		fmt.Fprintln(w, l.Get(ui.MsgStatsNoData))
		return
	}

	headers := []string{
		l.Get(ui.MsgHeaderMonth),
		l.Get(ui.MsgHeaderCount),
		l.Get(ui.MsgHeaderEstimate),
		l.Get(ui.MsgHeaderActual),
		l.Get(ui.MsgStatsAccuracy),
	}

	rows := make([][]string, len(stats))
	for i, s := range stats {
		rows[i] = []string{
			s.Month,
			fmt.Sprintf("%d", s.Count),
			ui.FormatHours(s.TotalEstimate),
			ui.FormatHours(s.TotalActual),
			formatAccuracy(s.TotalEstimate, s.TotalActual),
		}
	}

	ui.RenderTable(w, headers, rows)
}

// calculateAccuracy calculates estimate accuracy as a percentage
func calculateAccuracy(estimate, actual float64) float64 {
	if actual == 0 {
		return 0
	}
	return (estimate / actual) * 100
}

// formatAccuracy formats the accuracy as a percentage string
func formatAccuracy(estimate, actual float64) string {
	accuracy := calculateAccuracy(estimate, actual)
	return fmt.Sprintf("%.0f%%", accuracy)
}

// calculateAIRate calculates AI usage rate as a percentage
func calculateAIRate(aiMinutes int, actualHours float64) float64 {
	if actualHours == 0 {
		return 0
	}
	aiHours := float64(aiMinutes) / 60.0
	return (aiHours / actualHours) * 100
}

// formatAIRate formats the AI usage rate as a percentage string
func formatAIRate(aiMinutes int, actualHours float64) string {
	rate := calculateAIRate(aiMinutes, actualHours)
	return fmt.Sprintf("%.0f%%", rate)
}

// validateMonthFormat validates the month format (YYYY-MM)
func validateMonthFormat(month string) error {
	if month == "" {
		return ErrInvalidMonthFormat
	}

	var year, mon int
	n, err := fmt.Sscanf(month, "%d-%d", &year, &mon)
	if err != nil || n != 2 {
		return ErrInvalidMonthFormat
	}

	if mon < 1 || mon > 12 {
		return ErrInvalidMonthFormat
	}

	return nil
}
