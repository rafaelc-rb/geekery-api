package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/services"
)

type UserItemHandler struct {
	service *services.UserItemService
}

// NewUserItemHandler cria uma nova instância do handler de user items
func NewUserItemHandler(service *services.UserItemService) *UserItemHandler {
	return &UserItemHandler{service: service}
}

// getMockUserID retorna o userID mock para desenvolvimento
func getMockUserID() uint {
	// TODO: Substituir por autenticação real via JWT
	return 1
}

// AddToList adiciona um item à lista do usuário
// POST /api/my-list
func (h *UserItemHandler) AddToList(c *gin.Context) {
	ctx := c.Request.Context()
	userID := getMockUserID()

	var input struct {
		ItemID uint               `json:"item_id" binding:"required"`
		Status models.MediaStatus `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Status padrão se não fornecido
	if input.Status == "" {
		input.Status = models.StatusPlanned
	}

	userItem, err := h.service.AddToList(ctx, userID, input.ItemID, input.Status)
	if err != nil {
		if err == models.ErrDuplicateEntry {
			c.JSON(http.StatusConflict, gin.H{"error": "item already in your list"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userItem)
}

// GetMyList retorna a lista completa do usuário
// GET /api/my-list
func (h *UserItemHandler) GetMyList(c *gin.Context) {
	ctx := c.Request.Context()
	userID := getMockUserID()

	// Parâmetros de filtro opcionais
	statusParam := c.Query("status")
	favoriteParam := c.Query("favorite")

	var userItems []models.UserItem
	var err error

	// Filtrar por favoritos
	if favoriteParam == "true" {
		userItems, err = h.service.GetMyFavorites(ctx, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, userItems)
		return
	}

	// Filtrar por status
	if statusParam != "" {
		status := models.MediaStatus(statusParam)
		userItems, err = h.service.GetMyListByStatus(ctx, userID, status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, userItems)
		return
	}

	// Retornar lista completa
	userItems, err = h.service.GetMyList(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userItems)
}

// GetMyListItem retorna um item específico da lista do usuário
// GET /api/my-list/:id
func (h *UserItemHandler) GetMyListItem(c *gin.Context) {
	ctx := c.Request.Context()
	userID := getMockUserID()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	userItem, err := h.service.GetMyListItem(ctx, uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userItem)
}

// UpdateListItem atualiza um item da lista do usuário
// PUT /api/my-list/:id
func (h *UserItemHandler) UpdateListItem(c *gin.Context) {
	ctx := c.Request.Context()
	userID := getMockUserID()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var input struct {
		Status          models.MediaStatus  `json:"status"`
		Rating          float64             `json:"rating"`
		Favorite        bool                `json:"favorite"`
		Notes           string              `json:"notes"`
		ProgressType    models.ProgressType `json:"progress_type"`
		ProgressData    models.JSONB        `json:"progress_data"`
		CompletionCount int                 `json:"completion_count"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := &models.UserItem{
		Status:          input.Status,
		Rating:          input.Rating,
		Favorite:        input.Favorite,
		Notes:           input.Notes,
		ProgressType:    input.ProgressType,
		ProgressData:    input.ProgressData,
		CompletionCount: input.CompletionCount,
	}

	userItem, err := h.service.UpdateListItem(ctx, uint(id), userID, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userItem)
}

// RemoveFromList remove um item da lista do usuário
// DELETE /api/my-list/:id
func (h *UserItemHandler) RemoveFromList(c *gin.Context) {
	ctx := c.Request.Context()
	userID := getMockUserID()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.service.RemoveFromList(ctx, uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetStatistics retorna estatísticas da lista do usuário
// GET /api/my-list/stats
func (h *UserItemHandler) GetStatistics(c *gin.Context) {
	ctx := c.Request.Context()
	userID := getMockUserID()

	stats, err := h.service.GetStatistics(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
