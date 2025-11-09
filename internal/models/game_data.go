package models

import "gorm.io/gorm"

// GameData contém dados específicos de games
type GameData struct {
	gorm.Model
	ItemID          uint   `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Platform        string `json:"platform"`                            // PC, PS5, Switch, Xbox, etc
	Developer       string `json:"developer"`
	AveragePlaytime int    `json:"average_playtime"` // Horas médias para completar
}

// TableName especifica o nome da tabela
func (GameData) TableName() string {
	return "game_details"
}
