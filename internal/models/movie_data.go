package models

import "gorm.io/gorm"

// MovieData contém dados específicos de filmes
type MovieData struct {
	gorm.Model
	ItemID   uint   `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Director string `json:"director"`
	Runtime  int    `json:"runtime"` // em minutos
}

// TableName especifica o nome da tabela
func (MovieData) TableName() string {
	return "movie_details"
}
