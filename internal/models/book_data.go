package models

import (
	"time"

	"gorm.io/gorm"
)

// BookData contém dados específicos de livros/comics/novels
// Format especifica o subtipo: "manga", "manhwa", "light_novel", "web_novel", ou vazio para books
type BookData struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	ItemID    uint           `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Author    string         `json:"author"`
	Volumes   int            `json:"volumes"`            // Total de volumes (comic/novel sempre tem, books às vezes)
	Chapters  int            `json:"chapters"`           // Total de capítulos
	Pages     int            `json:"pages"`              // Total de páginas
	Format    string         `json:"format,omitempty"`   // Subtipo: manga, manhwa, light_novel, web_novel, etc
	Publisher string         `json:"publisher,omitempty"` // Editora/publisher
}

// TableName especifica o nome da tabela
func (BookData) TableName() string {
	return "book_details"
}
