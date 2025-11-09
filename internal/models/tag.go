package models

import (
	"gorm.io/gorm"
)

// Tag representa uma tag/categoria que pode ser associada a múltiplos items
type Tag struct {
	gorm.Model
	Name  string `json:"name" gorm:"uniqueIndex;not null"`
	Items []Item `json:"items,omitempty" gorm:"many2many:item_tags;"`
}

