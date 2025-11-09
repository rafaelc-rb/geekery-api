package models

import (
	"time"

	"gorm.io/gorm"
)

// SeriesData contém dados específicos de séries de TV
type SeriesData struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	ItemID    uint           `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Seasons   int            `json:"seasons"`
	Episodes  int            `json:"episodes"` // Total de episódios
}

// TableName especifica o nome da tabela
func (SeriesData) TableName() string {
	return "series_details"
}
