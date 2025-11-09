package models

import "gorm.io/gorm"

// SeriesData contém dados específicos de séries de TV
type SeriesData struct {
	gorm.Model
	ItemID   uint `json:"item_id" gorm:"uniqueIndex;not null"` // FK para items
	Seasons  int  `json:"seasons"`
	Episodes int  `json:"episodes"` // Total de episódios
}

// TableName especifica o nome da tabela
func (SeriesData) TableName() string {
	return "series_details"
}
