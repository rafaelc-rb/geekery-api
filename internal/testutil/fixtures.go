package testutil

import (
	"time"

	"github.com/rafaelc-rb/geekery-api/internal/models"
)

// ItemFixtures contém fixtures de items para testes
var ItemFixtures = struct {
	AnimeItem  models.Item
	MovieItem  models.Item
	SeriesItem models.Item
	GameItem   models.Item
	BookItem   models.Item
	ComicItem  models.Item
}{
	AnimeItem: models.Item{
		Title:       "Attack on Titan",
		Type:        models.MediaTypeAnime,
		Description: "Humans fight against titans",
		CoverURL:    "https://example.com/aot.jpg",
	},
	MovieItem: models.Item{
		Title:       "Inception",
		Type:        models.MediaTypeMovie,
		Description: "Dream within a dream",
		CoverURL:    "https://example.com/inception.jpg",
	},
	SeriesItem: models.Item{
		Title:       "Breaking Bad",
		Type:        models.MediaTypeSeries,
		Description: "Chemistry teacher turns meth cook",
		CoverURL:    "https://example.com/bb.jpg",
	},
	GameItem: models.Item{
		Title:       "The Witcher 3",
		Type:        models.MediaTypeGame,
		Description: "Epic RPG adventure",
		CoverURL:    "https://example.com/witcher.jpg",
	},
	BookItem: models.Item{
		Title:       "1984",
		Type:        models.MediaTypeBook,
		Description: "Dystopian masterpiece",
		CoverURL:    "https://example.com/1984.jpg",
	},
	ComicItem: models.Item{
		Title:       "One Piece",
		Type:        models.MediaTypeComic,
		Description: "Pirate adventure",
		CoverURL:    "https://example.com/onepiece.jpg",
	},
}

// TagFixtures contém fixtures de tags para testes
var TagFixtures = struct {
	ActionTag   models.Tag
	AdventureTag models.Tag
	SciFiTag    models.Tag
	DramaTag    models.Tag
}{
	ActionTag: models.Tag{
		Name: "action",
	},
	AdventureTag: models.Tag{
		Name: "adventure",
	},
	SciFiTag: models.Tag{
		Name: "sci-fi",
	},
	DramaTag: models.Tag{
		Name: "drama",
	},
}

// UserItemFixtures contém fixtures de user items para testes
var UserItemFixtures = struct {
	CompletedAnime models.UserItem
	InProgressGame models.UserItem
	PlannedMovie   models.UserItem
}{
	CompletedAnime: models.UserItem{
		UserID:          1,
		Status:          models.StatusCompleted,
		Rating:          9.5,
		Favorite:        true,
		Notes:           "Amazing anime!",
		ProgressType:    models.ProgressTypeEpisodic,
		CompletionCount: 1,
		ProgressData: models.JSONB{
			"episode": 75,
			"season":  4,
		},
	},
	InProgressGame: models.UserItem{
		UserID:          1,
		Status:          models.StatusInProgress,
		Rating:          8.5,
		Favorite:        false,
		Notes:           "Great game so far",
		ProgressType:    models.ProgressTypePercent,
		CompletionCount: 0,
		ProgressData: models.JSONB{
			"percent": 45,
			"hours":   30,
		},
	},
	PlannedMovie: models.UserItem{
		UserID:          1,
		Status:          models.StatusPlanned,
		Rating:          0,
		Favorite:        false,
		Notes:           "Want to watch",
		ProgressType:    models.ProgressTypeTime,
		CompletionCount: 0,
		ProgressData:    models.JSONB{},
	},
}

// NewTestItem cria um novo item para testes com valores padrão
func NewTestItem(title string, mediaType models.MediaType) *models.Item {
	return &models.Item{
		Title:       title,
		Type:        mediaType,
		Description: "Test description",
		ReleaseDate: timePtr(time.Now()),
	}
}

// NewTestTag cria uma nova tag para testes
func NewTestTag(name string) *models.Tag {
	return &models.Tag{
		Name: name,
	}
}

// NewTestUserItem cria um novo user item para testes
func NewTestUserItem(userID, itemID uint, status models.MediaStatus) *models.UserItem {
	return &models.UserItem{
		UserID:       userID,
		ItemID:       itemID,
		Status:       status,
		Rating:       0,
		Favorite:     false,
		ProgressType: models.ProgressTypeBoolean,
		ProgressData: models.JSONB{},
	}
}

// timePtr retorna um ponteiro para time.Time
func timePtr(t time.Time) *time.Time {
	return &t
}
