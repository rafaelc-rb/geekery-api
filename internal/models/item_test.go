package models

import (
	"strings"
	"testing"
	"time"
)

func TestItem_Validate(t *testing.T) {
	validDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	futureDate := time.Now().AddDate(100, 0, 0)
	pastDate := time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		item      *Item
		expectErr bool
		errorMsg  string
	}{
		{
			name: "valid_anime_item",
			item: &Item{
				Title:       "Attack on Titan",
				Type:        MediaTypeAnime,
				Description: "Anime about titans",
				ReleaseDate: &validDate,
			},
			expectErr: false,
		},
		{
			name: "missing_title",
			item: &Item{
				Type:        MediaTypeMovie,
				Description: "Movie without title",
			},
			expectErr: true,
			errorMsg:  "title",
		},
		{
			name: "title_too_long",
			item: &Item{
				Title: string(make([]byte, 501)),
				Type:  MediaTypeMovie,
			},
			expectErr: true,
			errorMsg:  "title",
		},
		{
			name: "invalid_media_type",
			item: &Item{
				Title: "Test",
				Type:  MediaType("invalid"),
			},
			expectErr: true,
			errorMsg:  "type",
		},
		{
			name: "release_date_too_far_future",
			item: &Item{
				Title:       "Future Movie",
				Type:        MediaTypeMovie,
				ReleaseDate: &futureDate,
			},
			expectErr: true,
			errorMsg:  "year",
		},
		{
			name: "release_date_too_old",
			item: &Item{
				Title:       "Ancient Movie",
				Type:        MediaTypeMovie,
				ReleaseDate: &pastDate,
			},
			expectErr: true,
			errorMsg:  "year",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if tt.expectErr && err != nil && tt.errorMsg != "" {
				if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.errorMsg)) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			}
		})
	}
}
