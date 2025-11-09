package models

import (
	"time"

	"gorm.io/gorm"
)

// MovieData contém dados específicos de filmes
type MovieData struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	ItemID    uint           `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Director  string         `json:"director"`
	Runtime   int            `json:"runtime"` // em minutos
}

// TableName especifica o nome da tabela
func (MovieData) TableName() string {
	return "movie_details"
}
