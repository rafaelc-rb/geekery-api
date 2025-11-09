package models

import (
	"testing"
)

func TestMediaType_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		mediaType MediaType
		want      bool
	}{
		{"anime_valid", MediaTypeAnime, true},
		{"movie_valid", MediaTypeMovie, true},
		{"series_valid", MediaTypeSeries, true},
		{"game_valid", MediaTypeGame, true},
		{"manga_valid", MediaTypeManga, true},
		{"light_novel_valid", MediaTypeLightNovel, true},
		{"music_valid", MediaTypeMusic, true},
		{"book_valid", MediaTypeBook, true},
		{"invalid_type", MediaType("invalid"), false},
		{"empty_type", MediaType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mediaType.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMediaStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status MediaStatus
		want   bool
	}{
		{"planned_valid", StatusPlanned, true},
		{"in_progress_valid", StatusInProgress, true},
		{"completed_valid", StatusCompleted, true},
		{"paused_valid", StatusPaused, true},
		{"dropped_valid", StatusDropped, true},
		{"invalid_status", MediaStatus("invalid"), false},
		{"empty_status", MediaStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgressType_IsValid(t *testing.T) {
	tests := []struct {
		name         string
		progressType ProgressType
		want         bool
	}{
		{"episodic_valid", ProgressTypeEpisodic, true},
		{"reading_valid", ProgressTypeReading, true},
		{"time_valid", ProgressTypeTime, true},
		{"percent_valid", ProgressTypePercent, true},
		{"boolean_valid", ProgressTypeBoolean, true},
		{"invalid_type", ProgressType("invalid"), false},
		{"empty_type", ProgressType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.progressType.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultProgressType(t *testing.T) {
	tests := []struct {
		name      string
		mediaType MediaType
		want      ProgressType
	}{
		{"anime_episodic", MediaTypeAnime, ProgressTypeEpisodic},
		{"series_episodic", MediaTypeSeries, ProgressTypeEpisodic},
		{"manga_reading", MediaTypeManga, ProgressTypeReading},
		{"light_novel_reading", MediaTypeLightNovel, ProgressTypeReading},
		{"book_reading", MediaTypeBook, ProgressTypeReading},
		{"movie_time", MediaTypeMovie, ProgressTypeTime},
		{"game_percent", MediaTypeGame, ProgressTypePercent},
		{"music_boolean", MediaTypeMusic, ProgressTypeBoolean},
		{"unknown_boolean", MediaType("unknown"), ProgressTypeBoolean},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDefaultProgressType(tt.mediaType)
			if got != tt.want {
				t.Errorf("GetDefaultProgressType() = %v, want %v", got, tt.want)
			}
		})
	}
}
