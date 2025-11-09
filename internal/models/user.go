package models

import (
	"gorm.io/gorm"
)

// User representa um usuário do sistema
type User struct {
	gorm.Model
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Name      string     `json:"name" gorm:"not null"`
	UserItems []UserItem `json:"user_items,omitempty" gorm:"foreignKey:UserID"`
}
