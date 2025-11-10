package data

import "github.com/rafaelc-rb/geekery-api/internal/models"

// GetTags retorna todas as tags para seed
func GetTags() []models.Tag {
	return []models.Tag{
		{Name: "Action"},
		{Name: "Adventure"},
		{Name: "Drama"},
		{Name: "Comedy"},
		{Name: "Fantasy"},
		{Name: "Sci-Fi"},
		{Name: "Shounen"},
		{Name: "Seinen"},
		{Name: "Romance"},
		{Name: "Mystery"},
		{Name: "Horror"},
		{Name: "Slice of Life"},
	}
}
