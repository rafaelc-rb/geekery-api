package models

import (
	"time"

	"gorm.io/gorm"
)

// MusicData contém dados específicos de música/álbuns
type MusicData struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	ItemID    uint           `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Artist    string         `json:"artist"`
	Album     string         `json:"album"`
	Duration  int            `json:"duration"` // Duração total do álbum em segundos
	Tracks    int            `json:"tracks"`   // Número de faixas
}

// TableName especifica o nome da tabela
func (MusicData) TableName() string {
	return "music_details"
}
