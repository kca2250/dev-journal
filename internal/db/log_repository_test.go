package db

import (
	"context"
	"testing"

	"github.com/kca2250/djou/internal/model"
	"github.com/kca2250/djou/internal/testutil"
)

func TestLogRepository_Create(t *testing.T) {
	sqlDB := testutil.SetupTestDB(t)
	repo := NewLogRepository(New(sqlDB))
	ctx := context.Background()

	tests := []struct {
		name    string
		input   model.LogInput
		wantErr bool
	}{
		{
			name: "valid input with all fields",
			input: model.LogInput{
				TaskName:      "テストタスク",
				EstimateHours: 2.0,
				ActualHours:   3.0,
				AIMinutes:     30,
				Problem:       "問題",
				Solution:      "解決",
				Learning:      "学び",
			},
			wantErr: false,
		},
		{
			name: "valid input with required fields only",
			input: model.LogInput{
				TaskName:      "必須のみタスク",
				EstimateHours: 1.0,
				ActualHours:   1.5,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, &tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogRepository_List(t *testing.T) {
	sqlDB := testutil.SetupTestDB(t)
	repo := NewLogRepository(New(sqlDB))
	ctx := context.Background()

	// テストデータを作成
	for i := 0; i < 15; i++ {
		input := model.LogInput{
			TaskName:      "タスク",
			EstimateHours: 1.0,
			ActualHours:   1.0,
		}
		if err := repo.Create(ctx, &input); err != nil {
			t.Fatalf("failed to create test data: %v", err)
		}
	}

	tests := []struct {
		name      string
		opts      ListOptions
		wantCount int
	}{
		{
			name:      "default limit 10",
			opts:      ListOptions{},
			wantCount: 10,
		},
		{
			name:      "custom limit 5",
			opts:      ListOptions{Limit: 5},
			wantCount: 5,
		},
		{
			name:      "limit exceeds total",
			opts:      ListOptions{Limit: 100},
			wantCount: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := repo.List(ctx, tt.opts)
			if err != nil {
				t.Errorf("List() error = %v", err)
				return
			}
			if len(logs) != tt.wantCount {
				t.Errorf("List() got %d logs, want %d", len(logs), tt.wantCount)
			}
		})
	}
}

func TestLogRepository_List_OrderByCreatedAtDesc(t *testing.T) {
	sqlDB := testutil.SetupTestDB(t)
	repo := NewLogRepository(New(sqlDB))
	ctx := context.Background()

	// Create logs with different task names
	names := []string{"first", "second", "third"}
	for _, name := range names {
		input := model.LogInput{
			TaskName:      name,
			EstimateHours: 1.0,
			ActualHours:   1.0,
		}
		if err := repo.Create(ctx, &input); err != nil {
			t.Fatalf("failed to create test data: %v", err)
		}
	}

	logs, err := repo.List(ctx, ListOptions{Limit: 3})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// IDs should be in descending order (newest first)
	// Since we insert first, second, third in order, IDs are 1, 2, 3
	// And they should be returned as 3, 2, 1 (newest first)
	if logs[0].ID < logs[1].ID || logs[1].ID < logs[2].ID {
		t.Errorf("logs should be in descending ID order, got IDs: %d, %d, %d",
			logs[0].ID, logs[1].ID, logs[2].ID)
	}
}

func TestLogRepository_Search(t *testing.T) {
	sqlDB := testutil.SetupTestDB(t)
	repo := NewLogRepository(New(sqlDB))
	ctx := context.Background()

	// Create test data
	testData := []model.LogInput{
		{TaskName: "ログイン実装", EstimateHours: 2, ActualHours: 3, Problem: "CORSエラー"},
		{TaskName: "API実装", EstimateHours: 1, ActualHours: 1, Learning: "CORSの設定が重要"},
		{TaskName: "DB設計", EstimateHours: 1, ActualHours: 2, Problem: "正規化の判断"},
		{TaskName: "テスト作成", EstimateHours: 1, ActualHours: 1},
	}

	for _, input := range testData {
		if err := repo.Create(ctx, &input); err != nil {
			t.Fatalf("failed to create test data: %v", err)
		}
	}

	tests := []struct {
		name      string
		keywords  []string
		limit     int
		wantCount int
	}{
		{
			name:      "single keyword in task name",
			keywords:  []string{"ログイン"},
			limit:     0,
			wantCount: 1,
		},
		{
			name:      "single keyword in problem",
			keywords:  []string{"CORS"},
			limit:     0,
			wantCount: 2,
		},
		{
			name:      "multiple keywords AND search",
			keywords:  []string{"CORS", "ログイン"},
			limit:     0,
			wantCount: 1,
		},
		{
			name:      "no match",
			keywords:  []string{"存在しない"},
			limit:     0,
			wantCount: 0,
		},
		{
			name:      "with limit",
			keywords:  []string{"実装"},
			limit:     1,
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := repo.Search(ctx, tt.keywords, tt.limit)
			if err != nil {
				t.Errorf("Search() error = %v", err)
				return
			}
			if len(logs) != tt.wantCount {
				t.Errorf("Search() got %d logs, want %d", len(logs), tt.wantCount)
			}
		})
	}
}

func TestLogRepository_GetStats(t *testing.T) {
	sqlDB := testutil.SetupTestDB(t)
	repo := NewLogRepository(New(sqlDB))
	ctx := context.Background()

	// Create test data
	testData := []model.LogInput{
		{TaskName: "タスク1", EstimateHours: 2, ActualHours: 3, AIMinutes: 30},
		{TaskName: "タスク2", EstimateHours: 1, ActualHours: 1, AIMinutes: 20},
		{TaskName: "タスク3", EstimateHours: 3, ActualHours: 4, AIMinutes: 60},
	}

	for _, input := range testData {
		if err := repo.Create(ctx, &input); err != nil {
			t.Fatalf("failed to create test data: %v", err)
		}
	}

	stats, err := repo.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}

	if stats.Count != 3 {
		t.Errorf("Count = %d, want 3", stats.Count)
	}
	if stats.TotalEstimate != 6 {
		t.Errorf("TotalEstimate = %f, want 6", stats.TotalEstimate)
	}
	if stats.TotalActual != 8 {
		t.Errorf("TotalActual = %f, want 8", stats.TotalActual)
	}
	if stats.TotalAIMinutes != 110 {
		t.Errorf("TotalAIMinutes = %d, want 110", stats.TotalAIMinutes)
	}
}

func TestLogRepository_GetStats_Empty(t *testing.T) {
	sqlDB := testutil.SetupTestDB(t)
	repo := NewLogRepository(New(sqlDB))
	ctx := context.Background()

	stats, err := repo.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}

	if stats.Count != 0 {
		t.Errorf("Count = %d, want 0", stats.Count)
	}
}
