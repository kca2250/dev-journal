package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kca2250/djou/internal/db"
	"github.com/kca2250/djou/internal/model"
	"github.com/kca2250/djou/internal/ui"
	"github.com/spf13/cobra"
)

var (
	searchLimit  int
	searchDetail bool
	searchTag    string
)

var searchCmd = &cobra.Command{
	Use:   "search [keywords...]",
	Short: "キーワードで開発日誌を検索",
	Long:  "キーワードで開発日誌を検索します。複数キーワードはAND検索になります。",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 0, "検索結果の表示件数を制限")
	searchCmd.Flags().BoolVarP(&searchDetail, "detail", "d", false, "詳細表示モード")
	searchCmd.Flags().StringVar(&searchTag, "tag", "", "タグでフィルタリング")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
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

	// Perform search
	ctx := context.Background()
	searchOpts := db.SearchOptions{
		Keywords: args,
		Limit:    searchLimit,
		Tag:      searchTag,
	}
	logs, err := repo.SearchWithOptions(ctx, searchOpts)
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Check if no results
	if len(logs) == 0 {
		fmt.Println(localizer.Get(ui.MsgSearchNoResults))
		return nil
	}

	// Show result count
	fmt.Printf("🔍 %s\n\n", localizer.Getf(ui.MsgSearchResults, len(logs)))

	// Render search results
	if searchDetail {
		renderSearchDetail(os.Stdout, logs, localizer)
	} else {
		renderSearchTable(os.Stdout, logs, args, localizer)
	}

	return nil
}

// getMatchFields returns a list of fields that match the keywords
func getMatchFields(log model.Log, keywords []string, l *ui.Localizer) []string {
	var matches []string

	for _, keyword := range keywords {
		kw := strings.ToLower(keyword)

		if strings.Contains(strings.ToLower(log.TaskName), kw) {
			if !contains(matches, l.Get(ui.MsgMatchTaskName)) {
				matches = append(matches, l.Get(ui.MsgMatchTaskName))
			}
		}
		if log.Memo != nil && strings.Contains(strings.ToLower(*log.Memo), kw) {
			if !contains(matches, l.Get(ui.MsgMatchMemo)) {
				matches = append(matches, l.Get(ui.MsgMatchMemo))
			}
		}
	}

	return matches
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func renderSearchTable(w *os.File, logs []model.Log, keywords []string, l *ui.Localizer) {
	headers := []string{
		l.Get(ui.MsgHeaderDate),
		l.Get(ui.MsgHeaderTask),
		l.Get(ui.MsgHeaderMatchField),
	}

	rows := make([][]string, len(logs))
	for i, log := range logs {
		matchFields := getMatchFields(log, keywords, l)
		rows[i] = []string{
			log.CreatedAt.Format("2006-01-02"),
			ui.TruncateString(log.TaskName, 20),
			strings.Join(matchFields, ", "),
		}
	}

	ui.RenderTable(w, headers, rows)
}

func renderSearchDetail(w *os.File, logs []model.Log, l *ui.Localizer) {
	headers := []string{
		l.Get(ui.MsgHeaderID),
		l.Get(ui.MsgHeaderDate),
		l.Get(ui.MsgHeaderTask),
		l.Get(ui.MsgHeaderEstimate),
		l.Get(ui.MsgHeaderActual),
		l.Get(ui.MsgHeaderMemo),
		l.Get(ui.MsgHeaderTags),
	}

	rows := make([][]string, len(logs))
	for i, log := range logs {
		memo := "-"
		if log.Memo != nil && *log.Memo != "" {
			memo = ui.TruncateString(*log.Memo, 30)
		}
		tags := "-"
		if log.Tags != nil && *log.Tags != "" {
			tags = ui.TruncateString(*log.Tags, 20)
		}
		rows[i] = []string{
			fmt.Sprintf("%d", log.ID),
			log.CreatedAt.Format("2006-01-02"),
			ui.TruncateString(log.TaskName, 20),
			ui.FormatHours(log.EstimateHours),
			ui.FormatHours(log.ActualHours),
			memo,
			tags,
		}
	}

	ui.RenderTable(w, headers, rows)
}
