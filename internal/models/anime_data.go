package models

import (
	"time"

	"gorm.io/gorm"
)

// AnimeData contém dados específicos de animes
type AnimeData struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	ItemID    uint           `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Episodes  int            `json:"episodes"`
	Studio    string         `json:"studio"`
}

// TableName especifica o nome da tabela
func (AnimeData) TableName() string {
	return "anime_details"
}
