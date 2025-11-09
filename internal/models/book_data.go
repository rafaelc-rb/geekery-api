package models

import "gorm.io/gorm"

// BookData contém dados específicos de livros/mangás/light novels
type BookData struct {
	gorm.Model
	ItemID   uint   `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Author   string `json:"author"`
	Volumes  int    `json:"volumes"`  // Total de volumes (manga/LN sempre tem, books às vezes)
	Chapters int    `json:"chapters"` // Total de capítulos
	Pages    int    `json:"pages"`    // Total de páginas
}

// TableName especifica o nome da tabela
func (BookData) TableName() string {
	return "book_details"
}
