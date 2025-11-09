package models

import (
	"time"

	"gorm.io/gorm"
)

// BookData contém dados específicos de livros/mangás/light novels
type BookData struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	ItemID    uint           `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Author    string         `json:"author"`
	Volumes   int            `json:"volumes"`  // Total de volumes (manga/LN sempre tem, books às vezes)
	Chapters  int            `json:"chapters"` // Total de capítulos
	Pages     int            `json:"pages"`    // Total de páginas
}

// TableName especifica o nome da tabela
func (BookData) TableName() string {
	return "book_details"
}
