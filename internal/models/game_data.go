package models

import (
	"time"

	"gorm.io/gorm"
)

// GameData contém dados específicos de games
type GameData struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	ItemID          uint           `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Platform        string         `json:"platform"`                            // PC, PS5, Switch, Xbox, etc
	Developer       string         `json:"developer"`
	AveragePlaytime int            `json:"average_playtime"` // Horas médias para completar
}

// TableName especifica o nome da tabela
func (GameData) TableName() string {
	return "game_details"
}
