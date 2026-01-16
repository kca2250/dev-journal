package ui

// Message keys for internationalization
const (
	// Record command messages
	MsgRecordSuccess       = "record_success"
	MsgRecordCancelled     = "record_cancelled"
	MsgTaskNameLabel       = "task_name_label"
	MsgEstimateHoursLabel  = "estimate_hours_label"
	MsgActualHoursLabel    = "actual_hours_label"
	MsgAIMinutesLabel      = "ai_minutes_label"
	MsgProblemLabel        = "problem_label"
	MsgSolutionLabel       = "solution_label"
	MsgLearningLabel       = "learning_label"
	MsgTaskNameRequired    = "task_name_required"
	MsgInvalidNumber       = "invalid_number"
	MsgPositiveRequired    = "positive_required"
	MsgNonNegativeRequired = "non_negative_required"

	// List command messages
	MsgNoLogs     = "no_logs"
	MsgLogsHeader = "logs_header"

	// Search command messages
	MsgSearchNoResults  = "search_no_results"
	MsgSearchResults    = "search_results"
	MsgHeaderMatchField = "header_match_field"
	MsgMatchTaskName    = "match_task_name"
	MsgMatchProblem     = "match_problem"
	MsgMatchSolution    = "match_solution"
	MsgMatchLearning    = "match_learning"
	MsgSelectDetail     = "select_detail"
	MsgDetailHeader     = "detail_header"

	// Stats command messages
	MsgStatsHeader       = "stats_header"
	MsgStatsTotalLogs    = "stats_total_logs"
	MsgStatsPeriod       = "stats_period"
	MsgStatsTotalHours   = "stats_total_hours"
	MsgStatsAccuracy     = "stats_accuracy"
	MsgStatsTotalAI      = "stats_total_ai"
	MsgStatsMonthlyTable = "stats_monthly_table"
	MsgStatsNoData       = "stats_no_data"

	// Export command messages
	MsgExportSuccess  = "export_success"
	MsgExportNoLogs   = "export_no_logs"
	MsgExportCreated  = "export_created"
	MsgExportFileNote = "export_file_note"

	// Table headers
	MsgHeaderID       = "header_id"
	MsgHeaderDate     = "header_date"
	MsgHeaderTask     = "header_task"
	MsgHeaderEstimate = "header_estimate"
	MsgHeaderActual   = "header_actual"
	MsgHeaderAI       = "header_ai"
	MsgHeaderMonth    = "header_month"
	MsgHeaderCount    = "header_count"

	// Common messages
	MsgError       = "error"
	MsgConfirmExit = "confirm_exit"
)

// messages contains all translations
var messages = map[string]map[string]string{
	"ja": {
		// Record command messages
		MsgRecordSuccess:       "ログを記録しました",
		MsgRecordCancelled:     "記録をキャンセルしました",
		MsgTaskNameLabel:       "タスク名",
		MsgEstimateHoursLabel:  "見積時間 (h)",
		MsgActualHoursLabel:    "実績時間 (h)",
		MsgAIMinutesLabel:      "AI使用時間 (min)",
		MsgProblemLabel:        "課題・問題点",
		MsgSolutionLabel:       "解決策",
		MsgLearningLabel:       "学び",
		MsgTaskNameRequired:    "タスク名は必須です",
		MsgInvalidNumber:       "有効な数値を入力してください",
		MsgPositiveRequired:    "正の数を入力してください",
		MsgNonNegativeRequired: "0以上の数を入力してください",

		// List command messages
		MsgNoLogs:     "ログがありません",
		MsgLogsHeader: "開発ログ一覧",

		// Search command messages
		MsgSearchNoResults:  "検索結果がありません",
		MsgSearchResults:    "検索結果: %d件",
		MsgHeaderMatchField: "マッチ箇所",
		MsgMatchTaskName:    "タスク名",
		MsgMatchProblem:     "ハマったこと",
		MsgMatchSolution:    "解決方法",
		MsgMatchLearning:    "学び",
		MsgSelectDetail:     "詳細を表示する番号を選択 (0で終了):",
		MsgDetailHeader:     "詳細情報",

		// Stats command messages
		MsgStatsHeader:       "統計情報",
		MsgStatsTotalLogs:    "総ログ数: %d件",
		MsgStatsPeriod:       "期間: %s 〜 %s",
		MsgStatsTotalHours:   "合計時間: 見積 %.1fh / 実績 %.1fh",
		MsgStatsAccuracy:     "見積精度: %.0f%%",
		MsgStatsTotalAI:      "AI使用時間: %dmin",
		MsgStatsMonthlyTable: "月別統計",
		MsgStatsNoData:       "データがありません",

		// Export command messages
		MsgExportSuccess:  "エクスポートが完了しました",
		MsgExportNoLogs:   "エクスポートするログがありません",
		MsgExportCreated:  "ファイルを作成しました: %s",
		MsgExportFileNote: "※同名ファイルが存在したため連番を追加しました",

		// Table headers
		MsgHeaderID:       "ID",
		MsgHeaderDate:     "日付",
		MsgHeaderTask:     "タスク",
		MsgHeaderEstimate: "見積",
		MsgHeaderActual:   "実績",
		MsgHeaderAI:       "AI",
		MsgHeaderMonth:    "月",
		MsgHeaderCount:    "件数",

		// Common messages
		MsgError:       "エラー: %s",
		MsgConfirmExit: "終了しますか？ (y/N)",
	},
	"en": {
		// Record command messages
		MsgRecordSuccess:       "Log recorded successfully",
		MsgRecordCancelled:     "Record cancelled",
		MsgTaskNameLabel:       "Task name",
		MsgEstimateHoursLabel:  "Estimate hours (h)",
		MsgActualHoursLabel:    "Actual hours (h)",
		MsgAIMinutesLabel:      "AI usage time (min)",
		MsgProblemLabel:        "Problems/Issues",
		MsgSolutionLabel:       "Solution",
		MsgLearningLabel:       "Learning",
		MsgTaskNameRequired:    "Task name is required",
		MsgInvalidNumber:       "Please enter a valid number",
		MsgPositiveRequired:    "Please enter a positive number",
		MsgNonNegativeRequired: "Please enter a non-negative number",

		// List command messages
		MsgNoLogs:     "No logs found",
		MsgLogsHeader: "Development Logs",

		// Search command messages
		MsgSearchNoResults:  "No results found",
		MsgSearchResults:    "Search results: %d",
		MsgHeaderMatchField: "Match Field",
		MsgMatchTaskName:    "Task Name",
		MsgMatchProblem:     "Problem",
		MsgMatchSolution:    "Solution",
		MsgMatchLearning:    "Learning",
		MsgSelectDetail:     "Select number to view details (0 to exit):",
		MsgDetailHeader:     "Details",

		// Stats command messages
		MsgStatsHeader:       "Statistics",
		MsgStatsTotalLogs:    "Total logs: %d",
		MsgStatsPeriod:       "Period: %s - %s",
		MsgStatsTotalHours:   "Total hours: Est %.1fh / Act %.1fh",
		MsgStatsAccuracy:     "Estimate accuracy: %.0f%%",
		MsgStatsTotalAI:      "AI usage time: %dmin",
		MsgStatsMonthlyTable: "Monthly Statistics",
		MsgStatsNoData:       "No data available",

		// Export command messages
		MsgExportSuccess:  "Export completed",
		MsgExportNoLogs:   "No logs to export",
		MsgExportCreated:  "File created: %s",
		MsgExportFileNote: "* Sequential number added due to existing file",

		// Table headers
		MsgHeaderID:       "ID",
		MsgHeaderDate:     "Date",
		MsgHeaderTask:     "Task",
		MsgHeaderEstimate: "Est",
		MsgHeaderActual:   "Act",
		MsgHeaderAI:       "AI",
		MsgHeaderMonth:    "Month",
		MsgHeaderCount:    "Count",

		// Common messages
		MsgError:       "Error: %s",
		MsgConfirmExit: "Exit? (y/N)",
	},
}
