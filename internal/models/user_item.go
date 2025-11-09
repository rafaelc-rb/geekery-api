package models

import (
	"gorm.io/gorm"
)

// UserItem representa um item na lista pessoal do usuário
// Relaciona um usuário com um item do catálogo, incluindo dados pessoais
// como status, rating, notas, etc.
type UserItem struct {
	gorm.Model
	UserID      uint        `json:"user_id" gorm:"not null;index:idx_user_item,unique"`
	ItemID      uint        `json:"item_id" gorm:"not null;index:idx_user_item,unique"`
	Status      MediaStatus `json:"status" gorm:"type:varchar(50);default:'planned';check:status IN ('planned','in_progress','completed','paused','dropped')"`
	Rating      float64     `json:"rating" gorm:"default:0"`
	Favorite    bool        `json:"favorite" gorm:"default:false"`
	Notes       string      `json:"notes" gorm:"type:text"`

	// Sistema de progresso flexível
	ProgressType ProgressType `json:"progress_type" gorm:"type:varchar(50);check:progress_type IN ('episodic','reading','time','percent','boolean')"`
	ProgressData JSONB        `json:"progress_data" gorm:"type:jsonb"` // Dados flexíveis de progresso + history

	// Contador de conclusões (quantas vezes completou 100%)
	CompletionCount int `json:"completion_count" gorm:"default:0"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Item Item `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// TableName especifica o nome da tabela no banco de dados
func (UserItem) TableName() string {
	return "user_items"
}

// Validate valida os dados do UserItem
func (ui *UserItem) Validate() error {
	if ui.UserID == 0 {
		return ErrUserIDRequired
	}

	if ui.ItemID == 0 {
		return ErrItemIDRequired
	}

	if ui.Status != "" && !ui.Status.IsValid() {
		return ErrInvalidStatus
	}

	if ui.Rating < 0 || ui.Rating > 10 {
		return ErrInvalidRating
	}

	if ui.ProgressType != "" && !ui.ProgressType.IsValid() {
		return ErrInvalidProgressType
	}

	if ui.CompletionCount < 0 {
		return ErrInvalidCompletionCount
	}

	return nil
}

// IsCompleted verifica se o item foi marcado como completo
func (ui *UserItem) IsCompleted() bool {
	return ui.Status == StatusCompleted
}

// IsFavorite verifica se o item é favorito
func (ui *UserItem) IsFavorite() bool {
	return ui.Favorite
}
