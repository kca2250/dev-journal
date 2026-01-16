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

var searchLimit int

var searchCmd = &cobra.Command{
	Use:   "search [keywords...]",
	Short: "キーワードで開発日誌を検索",
	Long:  "キーワードで開発日誌を検索します。複数キーワードはAND検索になります。",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 0, "検索結果の表示件数を制限")
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
	logs, err := repo.Search(ctx, args, searchLimit)
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

	// Render search results table
	renderSearchTable(os.Stdout, logs, args, localizer)

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
		if log.Problem != nil && strings.Contains(strings.ToLower(*log.Problem), kw) {
			if !contains(matches, l.Get(ui.MsgMatchProblem)) {
				matches = append(matches, l.Get(ui.MsgMatchProblem))
			}
		}
		if log.Solution != nil && strings.Contains(strings.ToLower(*log.Solution), kw) {
			if !contains(matches, l.Get(ui.MsgMatchSolution)) {
				matches = append(matches, l.Get(ui.MsgMatchSolution))
			}
		}
		if log.Learning != nil && strings.Contains(strings.ToLower(*log.Learning), kw) {
			if !contains(matches, l.Get(ui.MsgMatchLearning)) {
				matches = append(matches, l.Get(ui.MsgMatchLearning))
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
