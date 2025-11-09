package models

import "gorm.io/gorm"

// MusicData contém dados específicos de música/álbuns
type MusicData struct {
	gorm.Model
	ItemID   uint   `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Duration int    `json:"duration"` // Duração total do álbum em segundos
	Tracks   int    `json:"tracks"`   // Número de faixas
}

// TableName especifica o nome da tabela
func (MusicData) TableName() string {
	return "music_details"
}
