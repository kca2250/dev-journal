package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kca2250/djou/internal/db"
	"github.com/kca2250/djou/internal/model"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerRecordTool registers the djou_record tool
func registerRecordTool(s *server.MCPServer) {
	tool := mcp.NewTool("djou_record",
		mcp.WithDescription("開発日誌を記録する"),
		mcp.WithString("task_name",
			mcp.Required(),
			mcp.Description("タスク名"),
		),
		mcp.WithNumber("estimate_hours",
			mcp.Required(),
			mcp.Description("見積もり時間（時間）"),
		),
		mcp.WithNumber("actual_hours",
			mcp.Required(),
			mcp.Description("実績時間（時間）"),
		),
		mcp.WithString("memo",
			mcp.Description("メモ（課題、解決策、学びなど）"),
		),
		mcp.WithString("tags",
			mcp.Description("タグ（カンマ区切り）"),
		),
	)

	s.AddTool(tool, handleRecord)
}

func handleRecord(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments
	taskName, err := request.RequireString("task_name")
	if err != nil {
		return mcp.NewToolResultError("task_name は必須です"), nil
	}

	estimateHours, err := request.RequireFloat("estimate_hours")
	if err != nil {
		return mcp.NewToolResultError("estimate_hours は必須です"), nil
	}

	actualHours, err := request.RequireFloat("actual_hours")
	if err != nil {
		return mcp.NewToolResultError("actual_hours は必須です"), nil
	}

	memo := request.GetString("memo", "")
	tags := request.GetString("tags", "")

	// Open database
	database, err := db.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("データベース接続エラー: %s", err.Error())), nil
	}
	defer database.Close()

	// Create log input
	input := model.LogInput{
		TaskName:      taskName,
		EstimateHours: estimateHours,
		ActualHours:   actualHours,
		Memo:          memo,
		Tags:          tags,
	}

	// Save to database
	repo := db.NewLogRepository(database)
	if err := repo.Create(ctx, &input); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("保存エラー: %s", err.Error())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ 記録しました: %s", taskName)), nil
}

// registerListTool registers the djou_list tool
func registerListTool(s *server.MCPServer) {
	tool := mcp.NewTool("djou_list",
		mcp.WithDescription("記録一覧を取得する"),
		mcp.WithString("filter",
			mcp.Description("\"week\" または \"month\"、指定なしで直近"),
		),
		mcp.WithNumber("limit",
			mcp.Description("表示件数、デフォルト: 10"),
		),
	)

	s.AddTool(tool, handleList)
}

func handleList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filter := request.GetString("filter", "")
	limit := request.GetInt("limit", 10)

	// Open database
	database, err := db.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("データベース接続エラー: %s", err.Error())), nil
	}
	defer database.Close()

	repo := db.NewLogRepository(database)

	opts := db.ListOptions{
		Limit: limit,
	}
	switch filter {
	case "week":
		opts.Week = true
	case "month":
		opts.Month = "current"
	}

	logs, err := repo.List(ctx, opts)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("取得エラー: %s", err.Error())), nil
	}

	if len(logs) == 0 {
		return mcp.NewToolResultText("ログがありません"), nil
	}

	// Build markdown table
	var sb strings.Builder
	sb.WriteString("| 日付 | タスク | 見積もり | 実績 | メモ |\n")
	sb.WriteString("|---|---|---|---|---|\n")

	for _, log := range logs {
		memo := "-"
		if log.Memo != nil && *log.Memo != "" {
			memo = truncateString(*log.Memo, 20)
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %.1fh | %.1fh | %s |\n",
			log.CreatedAt.Format("2006-01-02"),
			log.TaskName,
			log.EstimateHours,
			log.ActualHours,
			memo,
		))
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}

// registerSearchTool registers the djou_search tool
func registerSearchTool(s *server.MCPServer) {
	tool := mcp.NewTool("djou_search",
		mcp.WithDescription("キーワードで検索する"),
		mcp.WithString("keywords",
			mcp.Required(),
			mcp.Description("検索キーワード（スペース区切りでAND検索）"),
		),
		mcp.WithNumber("limit",
			mcp.Description("表示件数、デフォルト: 無制限"),
		),
	)

	s.AddTool(tool, handleSearch)
}

func handleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keywordsStr, err := request.RequireString("keywords")
	if err != nil {
		return mcp.NewToolResultError("keywords は必須です"), nil
	}

	limit := request.GetInt("limit", 0)

	// Split keywords
	keywords := strings.Fields(keywordsStr)

	// Open database
	database, err := db.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("データベース接続エラー: %s", err.Error())), nil
	}
	defer database.Close()

	repo := db.NewLogRepository(database)
	logs, err := repo.Search(ctx, keywords, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("検索エラー: %s", err.Error())), nil
	}

	if len(logs) == 0 {
		return mcp.NewToolResultText("🔍 検索結果: 0件"), nil
	}

	// Build result
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔍 検索結果: %d件\n\n", len(logs)))
	sb.WriteString("| 日付 | タスク | マッチ箇所 |\n")
	sb.WriteString("|---|---|---|\n")

	for _, log := range logs {
		matchFields := getMatchFields(log, keywords)
		sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
			log.CreatedAt.Format("2006-01-02"),
			log.TaskName,
			strings.Join(matchFields, ", "),
		))
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func getMatchFields(log model.Log, keywords []string) []string {
	var matches []string

	for _, kw := range keywords {
		kwLower := strings.ToLower(kw)

		if strings.Contains(strings.ToLower(log.TaskName), kwLower) {
			if !containsString(matches, "タスク名") {
				matches = append(matches, "タスク名")
			}
		}
		if log.Memo != nil && strings.Contains(strings.ToLower(*log.Memo), kwLower) {
			if !containsString(matches, "メモ") {
				matches = append(matches, "メモ")
			}
		}
	}

	return matches
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// registerStatsTool registers the djou_stats tool
func registerStatsTool(s *server.MCPServer) {
	tool := mcp.NewTool("djou_stats",
		mcp.WithDescription("統計情報を取得する"),
		mcp.WithString("month",
			mcp.Description("特定月を指定（例: \"2025-01\"）、指定なしで全体統計"),
		),
		mcp.WithBoolean("monthly_list",
			mcp.Description("true で月別一覧を表示"),
		),
	)

	s.AddTool(tool, handleStats)
}

func handleStats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	month := request.GetString("month", "")
	monthlyList := request.GetBool("monthly_list", false)

	// Open database
	database, err := db.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("データベース接続エラー: %s", err.Error())), nil
	}
	defer database.Close()

	repo := db.NewLogRepository(database)

	var sb strings.Builder
	sb.WriteString("📊 Dev Journal 統計\n\n")

	if month != "" {
		// Specific month stats
		stats, err := repo.GetStatsByMonth(ctx, month)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("統計取得エラー: %s", err.Error())), nil
		}

		sb.WriteString(fmt.Sprintf("期間: %s\n", month))
		sb.WriteString(fmt.Sprintf("記録件数: %d件\n\n", stats.Count))
		sb.WriteString("⏱️ 作業時間\n")
		sb.WriteString(fmt.Sprintf("  見積もり合計: %.1fh\n", stats.TotalEstimate))
		sb.WriteString(fmt.Sprintf("  実績合計: %.1fh\n", stats.TotalActual))

		if stats.TotalActual > 0 {
			accuracy := (stats.TotalEstimate / stats.TotalActual) * 100
			sb.WriteString(fmt.Sprintf("  見積もり精度: %.0f%%\n", accuracy))
		}
	} else {
		// Overall stats
		stats, err := repo.GetStats(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("統計取得エラー: %s", err.Error())), nil
		}

		if stats.MinDate != nil && stats.MaxDate != nil {
			sb.WriteString(fmt.Sprintf("期間: %s 〜 %s\n",
				stats.MinDate.Format("2006-01-02"),
				stats.MaxDate.Format("2006-01-02")))
		}
		sb.WriteString(fmt.Sprintf("記録件数: %d件\n\n", stats.Count))
		sb.WriteString("⏱️ 作業時間\n")
		sb.WriteString(fmt.Sprintf("  見積もり合計: %.1fh\n", stats.TotalEstimate))
		sb.WriteString(fmt.Sprintf("  実績合計: %.1fh\n", stats.TotalActual))

		if stats.TotalActual > 0 {
			accuracy := (stats.TotalEstimate / stats.TotalActual) * 100
			sb.WriteString(fmt.Sprintf("  見積もり精度: %.0f%%\n", accuracy))
		}
	}

	if monthlyList {
		monthlyStats, err := repo.GetMonthlyStats(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("月別統計取得エラー: %s", err.Error())), nil
		}

		if len(monthlyStats) > 0 {
			sb.WriteString("\n📅 月別統計\n\n")
			sb.WriteString("| 月 | 件数 | 見積もり | 実績 | 精度 |\n")
			sb.WriteString("|---|---|---|---|---|\n")

			for _, ms := range monthlyStats {
				accuracy := ""
				if ms.TotalActual > 0 {
					accuracy = fmt.Sprintf("%.0f%%", (ms.TotalEstimate/ms.TotalActual)*100)
				}
				sb.WriteString(fmt.Sprintf("| %s | %d | %.1fh | %.1fh | %s |\n",
					ms.Month, ms.Count, ms.TotalEstimate, ms.TotalActual, accuracy))
			}
		}
	}

	return mcp.NewToolResultText(sb.String()), nil
}

// registerExportTool registers the djou_export tool
func registerExportTool(s *server.MCPServer) {
	tool := mcp.NewTool("djou_export",
		mcp.WithDescription("CSV出力する"),
		mcp.WithString("output",
			mcp.Description("出力先ディレクトリ、デフォルト: カレントディレクトリ"),
		),
		mcp.WithString("from",
			mcp.Description("期間開始日（YYYY-MM-DD）"),
		),
		mcp.WithString("to",
			mcp.Description("期間終了日（YYYY-MM-DD）"),
		),
	)

	s.AddTool(tool, handleExport)
}

func handleExport(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	outputDir := request.GetString("output", ".")
	fromStr := request.GetString("from", "")
	toStr := request.GetString("to", "")

	var fromDate, toDate *time.Time
	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			return mcp.NewToolResultError("from の形式が不正です（YYYY-MM-DD）"), nil
		}
		fromDate = &t
	}

	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			return mcp.NewToolResultError("to の形式が不正です（YYYY-MM-DD）"), nil
		}
		endOfDay := t.Add(24*time.Hour - time.Second)
		toDate = &endOfDay
	}

	// Validate output directory
	info, err := os.Stat(outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return mcp.NewToolResultError("出力先ディレクトリが存在しません"), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("ディレクトリエラー: %s", err.Error())), nil
	}
	if !info.IsDir() {
		return mcp.NewToolResultError("出力先がディレクトリではありません"), nil
	}

	// Open database
	database, err := db.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("データベース接続エラー: %s", err.Error())), nil
	}
	defer database.Close()

	repo := db.NewLogRepository(database)
	logs, err := repo.Export(ctx, fromDate, toDate)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("エクスポートエラー: %s", err.Error())), nil
	}

	if len(logs) == 0 {
		return mcp.NewToolResultText("エクスポートするログがありません"), nil
	}

	// Generate filename
	currentMonth := time.Now().Format("2006-01")
	filePath := generateFilename(outputDir, currentMonth)

	// Write CSV
	if err := writeCSV(filePath, logs); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("ファイル書き込みエラー: %s", err.Error())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ %s に出力しました", filepath.Base(filePath))), nil
}

func generateFilename(dir, month string) string {
	baseName := fmt.Sprintf("djou_%s", month)
	ext := ".csv"

	filePath := filepath.Join(dir, baseName+ext)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return filePath
	}

	for i := 1; ; i++ {
		filePath = filepath.Join(dir, fmt.Sprintf("%s_%d%s", baseName, i, ext))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return filePath
		}
	}
}

func writeCSV(filePath string, logs []model.Log) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write UTF-8 BOM for Excel compatibility
	if _, err := file.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return err
	}

	// Write header
	headers := "date,task_name,estimate_hours,actual_hours,memo,tags\n"
	if _, err := file.WriteString(headers); err != nil {
		return err
	}

	// Write data
	for _, log := range logs {
		memo := ""
		if log.Memo != nil {
			memo = escapeCSV(*log.Memo)
		}
		tags := ""
		if log.Tags != nil {
			tags = escapeCSV(*log.Tags)
		}

		row := fmt.Sprintf("%s,%s,%.1f,%.1f,%s,%s\n",
			log.CreatedAt.Format("2006-01-02"),
			escapeCSV(log.TaskName),
			log.EstimateHours,
			log.ActualHours,
			memo,
			tags,
		)
		if _, err := file.WriteString(row); err != nil {
			return err
		}
	}

	return nil
}

func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}

// registerUpdateTool registers the djou_update tool
func registerUpdateTool(s *server.MCPServer) {
	tool := mcp.NewTool("djou_update",
		mcp.WithDescription("開発日誌を更新する"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("更新対象のログID"),
		),
		mcp.WithString("task_name",
			mcp.Description("タスク名（指定しない場合は現在の値を維持）"),
		),
		mcp.WithNumber("estimate_hours",
			mcp.Description("見積もり時間（指定しない場合は現在の値を維持）"),
		),
		mcp.WithNumber("actual_hours",
			mcp.Description("実績時間（指定しない場合は現在の値を維持）"),
		),
		mcp.WithString("memo",
			mcp.Description("メモ（指定しない場合は現在の値を維持）"),
		),
		mcp.WithString("tags",
			mcp.Description("タグ（指定しない場合は現在の値を維持）"),
		),
	)

	s.AddTool(tool, handleUpdate)
}

func handleUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse ID (required)
	id, err := request.RequireInt("id")
	if err != nil {
		return mcp.NewToolResultError("id は必須です"), nil
	}

	// Open database
	database, err := db.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("データベース接続エラー: %s", err.Error())), nil
	}
	defer database.Close()

	repo := db.NewLogRepository(database)

	// Get existing log
	existingLog, err := repo.GetById(ctx, int64(id))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("ログが見つかりません: %s", err.Error())), nil
	}

	// Build input with partial updates
	input := model.LogInput{
		TaskName:      existingLog.TaskName,
		EstimateHours: existingLog.EstimateHours,
		ActualHours:   existingLog.ActualHours,
	}
	if existingLog.Memo != nil {
		input.Memo = *existingLog.Memo
	}
	if existingLog.Tags != nil {
		input.Tags = *existingLog.Tags
	}

	// Override with provided values
	if taskName := request.GetString("task_name", ""); taskName != "" {
		input.TaskName = taskName
	}
	// Note: GetFloat returns 0 as default, so we can't distinguish between "not provided" and "provided as 0"
	// For hours, 0 is not a valid value anyway, so we treat it as "not provided"
	if estimateHours := request.GetFloat("estimate_hours", 0); estimateHours > 0 {
		input.EstimateHours = estimateHours
	}
	if actualHours := request.GetFloat("actual_hours", 0); actualHours > 0 {
		input.ActualHours = actualHours
	}
	// For memo and tags, empty string means "clear the field" vs not provided
	// We use a sentinel value approach: if user provides empty string, it will be set
	if memo := request.GetString("memo", "\x00"); memo != "\x00" {
		input.Memo = memo
	}
	if tags := request.GetString("tags", "\x00"); tags != "\x00" {
		input.Tags = tags
	}

	// Update
	if err := repo.Update(ctx, int64(id), &input); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("更新エラー: %s", err.Error())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ ID %d を更新しました", id)), nil
}

// registerDeleteTool registers the djou_delete tool
func registerDeleteTool(s *server.MCPServer) {
	tool := mcp.NewTool("djou_delete",
		mcp.WithDescription("開発日誌を削除する"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("削除対象のログID"),
		),
	)

	s.AddTool(tool, handleDelete)
}

func handleDelete(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse ID (required)
	id, err := request.RequireInt("id")
	if err != nil {
		return mcp.NewToolResultError("id は必須です"), nil
	}

	// Open database
	database, err := db.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("データベース接続エラー: %s", err.Error())), nil
	}
	defer database.Close()

	repo := db.NewLogRepository(database)

	// Delete
	if err := repo.Delete(ctx, int64(id)); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("削除エラー: %s", err.Error())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ ID %d を削除しました", id)), nil
}
