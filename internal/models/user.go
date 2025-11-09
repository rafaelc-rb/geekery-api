package models

import (
	"time"

	"gorm.io/gorm"
)

// User representa um usuário do sistema
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Name      string         `json:"name" gorm:"not null"`
	UserItems []UserItem     `json:"user_items,omitempty" gorm:"foreignKey:UserID"`
}
