package dto

// ErrorResponse representa uma resposta de erro padronizada
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// SuccessResponse representa uma resposta de sucesso padronizada
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationMeta representa metadados de paginação
type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalItems   int64 `json:"total_items"`
	ItemsPerPage int   `json:"items_per_page"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// PaginatedResponse representa uma resposta paginada
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationParams representa os parâmetros de entrada para paginação
type PaginationParams struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// Normalize normaliza os parâmetros de paginação aplicando valores padrão
func (p *PaginationParams) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 || p.Limit > 100 {
		p.Limit = 20 // default
	}
}

// GetOffset calcula o offset para a query SQL baseado em page e limit
func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// NewPaginatedResponse cria uma resposta paginada com metadata calculado
func NewPaginatedResponse(data interface{}, page, limit int, total int64) *PaginatedResponse {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return &PaginatedResponse{
		Data: data,
		Pagination: PaginationMeta{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   total,
			ItemsPerPage: limit,
			HasNext:      page < totalPages,
			HasPrevious:  page > 1,
		},
	}
}

// HealthResponse representa a resposta do health check
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}
