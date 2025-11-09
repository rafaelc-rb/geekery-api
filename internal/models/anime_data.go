package models

import "gorm.io/gorm"

// AnimeData contém dados específicos de animes
type AnimeData struct {
	gorm.Model
	ItemID   uint   `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Episodes int    `json:"episodes"`
	Studio   string `json:"studio"`
}

// TableName especifica o nome da tabela
func (AnimeData) TableName() string {
	return "anime_details"
}
