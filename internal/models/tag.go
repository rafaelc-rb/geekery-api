package models

import (
	"time"

	"gorm.io/gorm"
)

// Tag representa uma tag/categoria que pode ser associada a múltiplos items
type Tag struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	Items     []Item         `json:"items,omitempty" gorm:"many2many:item_tags;"`
}
