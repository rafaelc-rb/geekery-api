package dto

import "time"

// ItemDTO representa um Item para resposta da API
type ItemDTO struct {
	ID               uint                   `json:"id"`
	Title            string                 `json:"title"`
	Type             string                 `json:"type"`
	Description      string                 `json:"description"`
	ReleaseDate      *time.Time             `json:"release_date,omitempty"`
	CoverURL         string                 `json:"cover_url,omitempty"`
	ExternalMetadata map[string]interface{} `json:"external_metadata,omitempty"`
	Tags             []TagDTO               `json:"tags,omitempty"`
	SpecificData     interface{}            `json:"specific_data,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// CreateItemRequest representa o payload de criação de item
type CreateItemRequest struct {
	Title            string                 `json:"title" binding:"required"`
	Type             string                 `json:"type" binding:"required"`
	Description      string                 `json:"description"`
	ReleaseDate      *time.Time             `json:"release_date"`
	CoverURL         string                 `json:"cover_url"`
	ExternalMetadata map[string]interface{} `json:"external_metadata"`
	TagIDs           []uint                 `json:"tag_ids"`
	SpecificData     interface{}            `json:"specific_data"`
}

// UpdateItemRequest representa o payload de atualização de item
type UpdateItemRequest struct {
	Title            string                 `json:"title"`
	Type             string                 `json:"type"`
	Description      string                 `json:"description"`
	ReleaseDate      *time.Time             `json:"release_date"`
	CoverURL         string                 `json:"cover_url"`
	ExternalMetadata map[string]interface{} `json:"external_metadata"`
	TagIDs           []uint                 `json:"tag_ids"`
}

// AnimeDataDTO representa dados específicos de anime
type AnimeDataDTO struct {
	Episodes int    `json:"episodes"`
	Studio   string `json:"studio"`
}

// MovieDataDTO representa dados específicos de filme
type MovieDataDTO struct {
	Director string `json:"director"`
	Runtime  int    `json:"runtime"` // em minutos
}

// SeriesDataDTO representa dados específicos de série
type SeriesDataDTO struct {
	Seasons  int `json:"seasons"`
	Episodes int `json:"episodes"`
}

// GameDataDTO representa dados específicos de game
type GameDataDTO struct {
	Platform        string `json:"platform"`
	Developer       string `json:"developer"`
	AveragePlaytime int    `json:"average_playtime"`
}

// BookDataDTO representa dados específicos de livro/comic/novel
type BookDataDTO struct {
	Author    string `json:"author"`
	Volumes   int    `json:"volumes"`
	Chapters  int    `json:"chapters"`
	Pages     int    `json:"pages"`
	Format    string `json:"format,omitempty"`
	Publisher string `json:"publisher,omitempty"`
}
