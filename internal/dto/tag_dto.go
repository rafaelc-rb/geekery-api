package dto

import "time"

// TagDTO representa uma Tag para resposta da API
type TagDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTagRequest representa o payload de criação de tag
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
}

// UpdateTagRequest representa o payload de atualização de tag
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
}

// TagListResponse representa uma lista de tags
type TagListResponse struct {
	Tags  []TagDTO `json:"tags"`
	Total int      `json:"total"`
}
