package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/kca2250/djou/internal/db"
	"github.com/kca2250/djou/internal/model"
	"github.com/kca2250/djou/internal/ui"
	"github.com/spf13/cobra"
)

var (
	recordTask     string
	recordEstimate float64
	recordActual   float64
	recordMemo     string
	recordTag      string
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "開発日誌を記録",
	Long: `開発日誌を記録します。
フラグを指定するとクイック記録モードになります。
フラグを指定しない場合はインタラクティブモードで入力を受け付けます。

例:
  djou record                                 # インタラクティブモード
  djou record -t "タスク名" -e 2 -a 3          # クイック記録
  djou record -t "タスク名" -e 2 -a 3 -m "メモ" --tag "task_type:新機能"`,
	RunE: runRecordCmd,
}

func init() {
	recordCmd.Flags().StringVarP(&recordTask, "task", "t", "", "タスク名")
	recordCmd.Flags().Float64VarP(&recordEstimate, "estimate", "e", 0, "見積時間（時間）")
	recordCmd.Flags().Float64VarP(&recordActual, "actual", "a", 0, "実績時間（時間）")
	recordCmd.Flags().StringVarP(&recordMemo, "memo", "m", "", "メモ")
	recordCmd.Flags().StringVar(&recordTag, "tag", "", "タグ（カンマ区切り）")
	rootCmd.AddCommand(recordCmd)
}

func runRecordCmd(cmd *cobra.Command, args []string) error {
	// Initialize localizer
	localizer := ui.NewAutoLocalizer()
	ui.SetGlobalLocalizer(localizer)

	// Check if quick record mode (all required flags provided)
	quickMode := recordTask != "" && recordEstimate > 0 && recordActual > 0

	var input *model.LogInput
	var err error

	if quickMode {
		// Quick record mode
		input = &model.LogInput{
			TaskName:      recordTask,
			EstimateHours: recordEstimate,
			ActualHours:   recordActual,
			Memo:          recordMemo,
			Tags:          recordTag,
		}

		// Validate
		if err := model.ValidateInput(input); err != nil {
			return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
		}
	} else {
		// Interactive mode
		form := ui.NewRecordForm(localizer)
		input, err = form.Run()
		if err != nil {
			if errors.Is(err, ui.ErrFormCancelled) {
				fmt.Println(localizer.Get(ui.MsgRecordCancelled))
				return nil
			}
			return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
		}
	}

	// Open database
	database, err := db.Open()
	if err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}
	defer database.Close()

	// Create repository and save
	repo := db.NewLogRepository(database)
	ctx := context.Background()
	if err := repo.Create(ctx, input); err != nil {
		return fmt.Errorf("%s", localizer.Getf(ui.MsgError, err.Error()))
	}

	// Show success message
	fmt.Println("✅ " + localizer.Get(ui.MsgRecordSuccess))

	return nil
}
