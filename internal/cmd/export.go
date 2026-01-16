package cmd

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kca2250/djou/internal/db"
	"github.com/kca2250/djou/internal/model"
	"github.com/kca2250/djou/internal/ui"
	"github.com/spf13/cobra"
)

var (
	exportOutput string
	exportFrom   string
	exportTo     string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "CSVファイルに出力",
	Long:  "記録した開発日誌をCSVファイルに出力します。",
	RunE:  runExport,
}

func init() {
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", ".", "出力先ディレクトリ")
	exportCmd.Flags().StringVar(&exportFrom, "from", "", "出力期間の開始日 (YYYY-MM-DD)")
	exportCmd.Flags().StringVar(&exportTo, "to", "", "出力期間の終了日 (YYYY-MM-DD)")
	rootCmd.AddCommand(exportCmd)
}

// ErrInvalidDateFormat is returned when date format is invalid
var ErrInvalidDateFormat = errors.New("invalid date format: use YYYY-MM-DD")

// ErrDirectoryNotFound is returned when output directory doesn't exist
var ErrDirectoryNotFound = errors.New("output directory not found")

func runExport(cmd *cobra.Command, args []string) error {
	// Initialize localizer
	localizer := ui.NewAutoLocalizer()
	ui.SetGlobalLocalizer(localizer)

	// Validate output directory
	info, err := os.Stat(exportOutput)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrDirectoryNotFound
		}
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}
	if !info.IsDir() {
		return ErrDirectoryNotFound
	}

	// Parse date flags
	var fromDate, toDate *time.Time
	if exportFrom != "" {
		t, err := parseDateFlag(exportFrom)
		if err != nil {
			return err
		}
		fromDate = &t
	}
	if exportTo != "" {
		t, err := parseDateFlag(exportTo)
		if err != nil {
			return err
		}
		// Set to end of day
		endOfDay := t.Add(24*time.Hour - time.Second)
		toDate = &endOfDay
	}

	// Open database
	database, err := db.Open()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}
	defer database.Close()

	// Create repository
	repo := db.NewLogRepository(database)
	ctx := context.Background()

	// Fetch logs
	logs, err := repo.Export(ctx, fromDate, toDate)
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Check if no logs found
	if len(logs) == 0 {
		fmt.Println(localizer.Get(ui.MsgExportNoLogs))
		return nil
	}

	// Generate filename
	currentMonth := time.Now().Format("2006-01")
	filePath := generateFilename(exportOutput, currentMonth)

	// Write CSV
	if err := writeCSV(filePath, logs); err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Print success message
	fmt.Println(localizer.Get(ui.MsgExportSuccess))
	fmt.Println(localizer.Getf(ui.MsgExportCreated, filePath))

	return nil
}

// parseDateFlag parses a date string in YYYY-MM-DD format
func parseDateFlag(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, ErrInvalidDateFormat
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}

	return t, nil
}

// generateFilename generates a unique filename with sequential numbering
func generateFilename(dir, month string) string {
	baseName := fmt.Sprintf("djou_%s", month)
	ext := ".csv"

	// Check if base file exists
	filePath := filepath.Join(dir, baseName+ext)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return filePath
	}

	// Find next available number
	for i := 1; ; i++ {
		filePath = filepath.Join(dir, fmt.Sprintf("%s_%d%s", baseName, i, ext))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return filePath
		}
	}
}

// getCSVHeaders returns the CSV header row
func getCSVHeaders() []string {
	return []string{
		"date",
		"task_name",
		"estimate_hours",
		"actual_hours",
		"ai_minutes",
		"problem",
		"solution",
		"learning",
	}
}

// formatCSVRow formats a log entry as a CSV row
func formatCSVRow(log model.Log) []string {
	aiMinutes := ""
	if log.AIMinutes != nil {
		aiMinutes = fmt.Sprintf("%d", *log.AIMinutes)
	}

	problem := ""
	if log.Problem != nil {
		problem = *log.Problem
	}

	solution := ""
	if log.Solution != nil {
		solution = *log.Solution
	}

	learning := ""
	if log.Learning != nil {
		learning = *log.Learning
	}

	return []string{
		log.CreatedAt.Format("2006-01-02"),
		log.TaskName,
		fmt.Sprintf("%.1f", log.EstimateHours),
		fmt.Sprintf("%.1f", log.ActualHours),
		aiMinutes,
		problem,
		solution,
		learning,
	}
}

// writeCSV writes logs to a CSV file with UTF-8 BOM
func writeCSV(filePath string, logs []model.Log) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write UTF-8 BOM for Excel compatibility
	if _, err := file.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return fmt.Errorf("failed to write BOM: %w", err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(getCSVHeaders()); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// Write data rows
	for _, log := range logs {
		if err := writer.Write(formatCSVRow(log)); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}
