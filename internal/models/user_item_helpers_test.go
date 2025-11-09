package models

import (
	"testing"

	"gorm.io/gorm"
)

func TestUserItem_SetEpisodicProgress(t *testing.T) {
	tests := []struct {
		name    string
		season  int
		episode int
	}{
		{"set_season_1_episode_5", 1, 5},
		{"set_season_4_episode_12", 4, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserItem{}
			ui.SetEpisodicProgress(tt.season, tt.episode)

			if ui.ProgressType != ProgressTypeEpisodic {
				t.Errorf("Expected ProgressType %v, got %v", ProgressTypeEpisodic, ui.ProgressType)
			}
			if got := getInt(ui.ProgressData["season"]); got != tt.season {
				t.Errorf("Expected season %d, got %d", tt.season, got)
			}
			if got := getInt(ui.ProgressData["episode"]); got != tt.episode {
				t.Errorf("Expected episode %d, got %d", tt.episode, got)
			}
		})
	}
}

func TestUserItem_SetReadingProgress(t *testing.T) {
	tests := []struct {
		name    string
		chapter *int
		volume  *int
		page    *int
	}{
		{"only_chapter", intPtr(50), nil, nil},
		{"only_volume", nil, intPtr(5), nil},
		{"only_page", nil, nil, intPtr(150)},
		{"chapter_and_volume", intPtr(25), intPtr(3), nil},
		{"all_fields", intPtr(30), intPtr(4), intPtr(200)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserItem{}
			ui.SetReadingProgress(tt.chapter, tt.volume, tt.page)

			if ui.ProgressType != ProgressTypeReading {
				t.Errorf("Expected ProgressType %v, got %v", ProgressTypeReading, ui.ProgressType)
			}

			if tt.chapter != nil {
				if got := getInt(ui.ProgressData["chapter"]); got != *tt.chapter {
					t.Errorf("Expected chapter %d, got %d", *tt.chapter, got)
				}
			}
			if tt.volume != nil {
				if got := getInt(ui.ProgressData["volume"]); got != *tt.volume {
					t.Errorf("Expected volume %d, got %d", *tt.volume, got)
				}
			}
			if tt.page != nil {
				if got := getInt(ui.ProgressData["page"]); got != *tt.page {
					t.Errorf("Expected page %d, got %d", *tt.page, got)
				}
			}
		})
	}
}

func TestUserItem_StartNewView(t *testing.T) {
	ui := &UserItem{
		ProgressData: JSONB{},
		ProgressType: ProgressTypeEpisodic, // Definir tipo antes
	}

	ui.StartNewView()

	history := ui.getHistory()
	if len(history) != 1 {
		t.Errorf("Expected history length 1, got %d", len(history))
	}

	if !ui.IsCurrentViewInProgress() {
		t.Error("Expected current view to be in progress")
	}
	if got := ui.GetCurrentViewNumber(); got != 1 {
		t.Errorf("Expected current view number 1, got %d", got)
	}
}

func TestUserItem_CompleteCurrentView(t *testing.T) {
	ui := &UserItem{
		ProgressData: JSONB{},
		ProgressType: ProgressTypeEpisodic, // Definir tipo antes
	}

	// Iniciar uma view
	ui.StartNewView()
	if ui.CompletionCount != 0 {
		t.Errorf("Expected CompletionCount 0, got %d", ui.CompletionCount)
	}

	// Completar view
	ui.CompleteCurrentView()

	if ui.CompletionCount != 1 {
		t.Errorf("Expected CompletionCount 1, got %d", ui.CompletionCount)
	}
	if ui.IsCurrentViewInProgress() {
		t.Error("Expected current view to not be in progress")
	}

	history := ui.GetAllViews()
	if len(history) != 1 {
		t.Errorf("Expected history length 1, got %d", len(history))
	}

	if len(history) > 0 && history[0]["finished_at"] == nil {
		t.Error("Expected finished_at to be set")
	}
}

func TestUserItem_GetProgressPercent_Episodic(t *testing.T) {
	item := &Item{
		Model: gorm.Model{ID: 1}, // Adicionar ID
		Type:  MediaTypeAnime,
		AnimeData: &AnimeData{
			Episodes: 24,
		},
	}

	ui := &UserItem{
		Item: *item,
	}

	ui.SetEpisodicProgress(1, 12)

	percent := ui.GetProgressPercent()
	if percent != 50.0 {
		t.Errorf("Expected 50.0%%, got %.1f%%", percent)
	}
}

func TestUserItem_GetProgressPercent_Reading(t *testing.T) {
	tests := []struct {
		name            string
		chapters        int
		volumes         int
		pages           int
		currentChapter  *int
		currentVolume   *int
		currentPage     *int
		expectedPercent float64
	}{
		{"chapter_progress", 100, 0, 0, intPtr(50), nil, nil, 50.0},
		{"page_progress", 0, 0, 300, nil, nil, intPtr(150), 50.0},
		{"volume_progress", 0, 10, 0, nil, intPtr(5), nil, 50.0},
		{"chapter_priority", 100, 10, 300, intPtr(25), intPtr(2), intPtr(50), 25.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &Item{
				Model: gorm.Model{ID: 1}, // Adicionar ID
				Type:  MediaTypeManga,
				BookData: &BookData{
					Chapters: tt.chapters,
					Volumes:  tt.volumes,
					Pages:    tt.pages,
				},
			}

			ui := &UserItem{Item: *item}
			ui.SetReadingProgress(tt.currentChapter, tt.currentVolume, tt.currentPage)

			percent := ui.GetProgressPercent()
			if percent != tt.expectedPercent {
				t.Errorf("Expected %.1f%%, got %.1f%%", tt.expectedPercent, percent)
			}
		})
	}
}

func TestUserItem_Validate(t *testing.T) {
	tests := []struct {
		name      string
		userItem  *UserItem
		expectErr bool
	}{
		{
			name: "valid_user_item",
			userItem: &UserItem{
				UserID:       1,
				ItemID:       1,
				Status:       StatusCompleted,
				Rating:       8.5,
				ProgressType: ProgressTypeEpisodic,
			},
			expectErr: false,
		},
		{
			name: "missing_user_id",
			userItem: &UserItem{
				ItemID:       1,
				Status:       StatusCompleted,
				ProgressType: ProgressTypeEpisodic,
			},
			expectErr: true,
		},
		{
			name: "invalid_rating_too_high",
			userItem: &UserItem{
				UserID:       1,
				ItemID:       1,
				Status:       StatusCompleted,
				Rating:       11.0,
				ProgressType: ProgressTypeEpisodic,
			},
			expectErr: true,
		},
		{
			name: "invalid_status",
			userItem: &UserItem{
				UserID:       1,
				ItemID:       1,
				Status:       MediaStatus("invalid"),
				ProgressType: ProgressTypeEpisodic,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.userItem.Validate()

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}
