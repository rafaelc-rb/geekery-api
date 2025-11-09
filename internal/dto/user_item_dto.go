package dto

import "time"

// UserItemDTO representa um UserItem para resposta da API
type UserItemDTO struct {
	ID              uint                   `json:"id"`
	UserID          uint                   `json:"user_id"`
	ItemID          uint                   `json:"item_id"`
	Status          string                 `json:"status"`
	Rating          float64                `json:"rating"`
	Favorite        bool                   `json:"favorite"`
	Notes           string                 `json:"notes,omitempty"`
	ProgressType    string                 `json:"progress_type"`
	ProgressData    map[string]interface{} `json:"progress_data,omitempty"`
	CompletionCount int                    `json:"completion_count"`
	Item            *ItemDTO               `json:"item,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// AddToListRequest representa o payload para adicionar item à lista
type AddToListRequest struct {
	ItemID uint   `json:"item_id" binding:"required"`
	Status string `json:"status"` // opcional, padrão: "planned"
}

// UpdateUserItemRequest representa o payload de atualização de user item
type UpdateUserItemRequest struct {
	Status          string                 `json:"status"`
	Rating          float64                `json:"rating"`
	Favorite        bool                   `json:"favorite"`
	Notes           string                 `json:"notes"`
	ProgressType    string                 `json:"progress_type"`
	ProgressData    map[string]interface{} `json:"progress_data"`
	CompletionCount int                    `json:"completion_count"`
}

// UpdateProgressRequest representa payload específico para atualizar progresso
type UpdateProgressRequest struct {
	ProgressType string                 `json:"progress_type"`
	ProgressData map[string]interface{} `json:"progress_data" binding:"required"`
}

// UserListStatsDTO representa estatísticas da lista do usuário
type UserListStatsDTO struct {
	Total       int64            `json:"total"`
	InProgress  int64            `json:"in_progress"`
	Completed   int64            `json:"completed"`
	Planned     int64            `json:"planned"`
	Paused      int64            `json:"paused"`
	Dropped     int64            `json:"dropped"`
	Favorites   int64            `json:"favorites"`
	ByMediaType map[string]int64 `json:"by_media_type,omitempty"`
}
