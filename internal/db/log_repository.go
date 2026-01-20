package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kca2250/djou/internal/model"
)

// ListOptions represents options for listing logs
type ListOptions struct {
	Limit int
	Week  bool
	Month string // empty for current month, or "YYYY-MM" for specific month
	Tag   string // tag filter
}

// GetLimit returns the limit or default value
func (o ListOptions) GetLimit() int {
	if o.Limit <= 0 {
		return 10
	}
	return o.Limit
}

// Stats represents aggregated statistics
type Stats struct {
	Count         int
	MinDate       *time.Time
	MaxDate       *time.Time
	TotalEstimate float64
	TotalActual   float64
}

// MonthlyStats represents monthly aggregated statistics
type MonthlyStats struct {
	Month         string
	Count         int
	TotalEstimate float64
	TotalActual   float64
}

// LogRepository handles database operations for logs
type LogRepository struct {
	db *DB
}

// NewLogRepository creates a new LogRepository
func NewLogRepository(db *DB) *LogRepository {
	return &LogRepository{db: db}
}

// Create inserts a new log entry
func (r *LogRepository) Create(ctx context.Context, input *model.LogInput) error {
	query := `
		INSERT INTO logs (task_name, estimate_hours, actual_hours, memo, tags)
		VALUES (?, ?, ?, ?, ?)
	`

	var memo, tags *string
	if input.Memo != "" {
		memo = &input.Memo
	}
	if input.Tags != "" {
		tags = &input.Tags
	}

	_, err := r.db.ExecContext(ctx, query,
		input.TaskName,
		input.EstimateHours,
		input.ActualHours,
		memo,
		tags,
	)
	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	return nil
}

// List retrieves logs with optional filters
func (r *LogRepository) List(ctx context.Context, opts ListOptions) ([]model.Log, error) {
	query := `
		SELECT id, created_at, task_name, estimate_hours, actual_hours, memo, tags
		FROM logs
	`
	args := []any{}

	// Add WHERE clause for week/month filters
	var whereClauses []string

	if opts.Week {
		// Get start of current week (Monday)
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startOfWeek := now.AddDate(0, 0, -(weekday - 1))
		startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
		endOfWeek := startOfWeek.AddDate(0, 0, 7)

		whereClauses = append(whereClauses, "created_at >= ? AND created_at < ?")
		args = append(args, startOfWeek, endOfWeek)
	}

	if opts.Month != "" {
		// Parse month or use current month
		var year, month int
		if opts.Month == "current" {
			now := time.Now()
			year = now.Year()
			month = int(now.Month())
		} else {
			_, err := fmt.Sscanf(opts.Month, "%d-%d", &year, &month)
			if err != nil {
				return nil, fmt.Errorf("invalid month format: %w", err)
			}
		}

		startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		endOfMonth := startOfMonth.AddDate(0, 1, 0)

		whereClauses = append(whereClauses, "created_at >= ? AND created_at < ?")
		args = append(args, startOfMonth, endOfMonth)
	}

	// Tag filter
	if opts.Tag != "" {
		whereClauses = append(whereClauses, "tags LIKE ?")
		args = append(args, "%"+opts.Tag+"%")
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY created_at DESC, id DESC"
	query += fmt.Sprintf(" LIMIT %d", opts.GetLimit())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}
	defer rows.Close()

	return scanLogs(rows)
}

// Search searches logs by keywords (AND search)
func (r *LogRepository) Search(ctx context.Context, keywords []string, limit int) ([]model.Log, error) {
	if len(keywords) == 0 {
		return []model.Log{}, nil
	}

	query := `
		SELECT id, created_at, task_name, estimate_hours, actual_hours, memo, tags
		FROM logs
		WHERE 1=1
	`
	args := []any{}

	// AND search: all keywords must match in any field
	for _, keyword := range keywords {
		query += ` AND (
			task_name LIKE ? OR
			memo LIKE ?
		)`
		pattern := "%" + keyword + "%"
		args = append(args, pattern, pattern)
	}

	query += " ORDER BY created_at DESC, id DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search logs: %w", err)
	}
	defer rows.Close()

	return scanLogs(rows)
}

// GetStats returns aggregated statistics for all logs
func (r *LogRepository) GetStats(ctx context.Context) (*Stats, error) {
	query := `
		SELECT
			COUNT(*) as count,
			MIN(created_at) as min_date,
			MAX(created_at) as max_date,
			COALESCE(SUM(estimate_hours), 0) as total_estimate,
			COALESCE(SUM(actual_hours), 0) as total_actual
		FROM logs
	`

	var stats Stats
	var minDate, maxDate sql.NullString

	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.Count,
		&minDate,
		&maxDate,
		&stats.TotalEstimate,
		&stats.TotalActual,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	if minDate.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", minDate.String)
		if err == nil {
			stats.MinDate = &t
		}
	}
	if maxDate.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", maxDate.String)
		if err == nil {
			stats.MaxDate = &t
		}
	}

	return &stats, nil
}

// GetMonthlyStats returns monthly statistics
func (r *LogRepository) GetMonthlyStats(ctx context.Context) ([]MonthlyStats, error) {
	query := `
		SELECT
			strftime('%Y-%m', created_at) as month,
			COUNT(*) as count,
			SUM(estimate_hours) as total_estimate,
			SUM(actual_hours) as total_actual
		FROM logs
		GROUP BY strftime('%Y-%m', created_at)
		ORDER BY month DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly stats: %w", err)
	}
	defer rows.Close()

	var result []MonthlyStats
	for rows.Next() {
		var ms MonthlyStats
		if err := rows.Scan(&ms.Month, &ms.Count, &ms.TotalEstimate, &ms.TotalActual); err != nil {
			return nil, fmt.Errorf("failed to scan monthly stats: %w", err)
		}
		result = append(result, ms)
	}

	return result, rows.Err()
}

// GetStatsByMonth returns statistics for a specific month
func (r *LogRepository) GetStatsByMonth(ctx context.Context, month string) (*Stats, error) {
	var year, mon int
	_, err := fmt.Sscanf(month, "%d-%d", &year, &mon)
	if err != nil {
		return nil, fmt.Errorf("invalid month format: %w", err)
	}

	startOfMonth := time.Date(year, time.Month(mon), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	query := `
		SELECT
			COUNT(*) as count,
			MIN(created_at) as min_date,
			MAX(created_at) as max_date,
			COALESCE(SUM(estimate_hours), 0) as total_estimate,
			COALESCE(SUM(actual_hours), 0) as total_actual
		FROM logs
		WHERE created_at >= ? AND created_at < ?
	`

	var stats Stats
	var minDate, maxDate sql.NullString

	err = r.db.QueryRowContext(ctx, query, startOfMonth, endOfMonth).Scan(
		&stats.Count,
		&minDate,
		&maxDate,
		&stats.TotalEstimate,
		&stats.TotalActual,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats by month: %w", err)
	}

	if minDate.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", minDate.String)
		if err == nil {
			stats.MinDate = &t
		}
	}
	if maxDate.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", maxDate.String)
		if err == nil {
			stats.MaxDate = &t
		}
	}

	return &stats, nil
}

// Export retrieves logs for export with optional date range
func (r *LogRepository) Export(ctx context.Context, from, to *time.Time) ([]model.Log, error) {
	query := `
		SELECT id, created_at, task_name, estimate_hours, actual_hours, memo, tags
		FROM logs
	`
	args := []any{}
	var whereClauses []string

	if from != nil {
		whereClauses = append(whereClauses, "created_at >= ?")
		args = append(args, *from)
	}
	if to != nil {
		whereClauses = append(whereClauses, "created_at < ?")
		args = append(args, *to)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY created_at ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to export logs: %w", err)
	}
	defer rows.Close()

	return scanLogs(rows)
}

// Update updates an existing log entry
func (r *LogRepository) Update(ctx context.Context, id int64, input *model.LogInput) error {
	query := `
		UPDATE logs
		SET task_name = ?, estimate_hours = ?, actual_hours = ?, memo = ?, tags = ?
		WHERE id = ?
	`

	var memo, tags *string
	if input.Memo != "" {
		memo = &input.Memo
	}
	if input.Tags != "" {
		tags = &input.Tags
	}

	result, err := r.db.ExecContext(ctx, query,
		input.TaskName,
		input.EstimateHours,
		input.ActualHours,
		memo,
		tags,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("log not found: %d", id)
	}

	return nil
}

// Delete deletes a log entry by ID
func (r *LogRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM logs WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("log not found: %d", id)
	}

	return nil
}

// scanLogs scans rows into Log structs
func scanLogs(rows *sql.Rows) ([]model.Log, error) {
	var logs []model.Log

	for rows.Next() {
		var log model.Log
		err := rows.Scan(
			&log.ID,
			&log.CreatedAt,
			&log.TaskName,
			&log.EstimateHours,
			&log.ActualHours,
			&log.Memo,
			&log.Tags,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}
