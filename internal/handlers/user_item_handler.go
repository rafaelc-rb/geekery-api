package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/dto"
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
// @Summary      Add item to my list
// @Description  Add a catalog item to user's personal tracking list
// @Tags         my-list
// @Accept       json
// @Produce      json
// @Param        item  body  object{item_id=int,status=string}  true  "Item to add (item_id required, status optional)"
// @Success      201  {object}  models.UserItem       "Item added successfully"
// @Failure      400  {object}  map[string]string     "Bad request - validation error"
// @Failure      409  {object}  map[string]string     "Conflict - item already in list"
// @Router       /my-list [post]
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

// GetMyList retorna a lista completa do usuário com paginação
// @Summary      Get my list
// @Description  Get user's personal tracking list with pagination and filters
// @Tags         my-list
// @Accept       json
// @Produce      json
// @Param        page      query  int     false  "Page number" default(1)
// @Param        limit     query  int     false  "Items per page" default(20)
// @Param        status    query  string  false  "Filter by status" Enums(planned, in_progress, completed, paused, dropped)
// @Param        favorite  query  bool    false  "Filter favorites only"
// @Success      200  {object}  dto.PaginatedResponse  "Success - returns user's items"
// @Failure      400  {object}  map[string]string      "Bad request - invalid parameters"
// @Failure      500  {object}  map[string]string      "Internal server error"
// @Router       /my-list [get]
func (h *UserItemHandler) GetMyList(c *gin.Context) {
	ctx := c.Request.Context()
	userID := getMockUserID()

	// Parse parâmetros de paginação
	var params dto.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}
	params.Normalize()

	// Parâmetros de filtro opcionais
	statusParam := c.Query("status")
	favoriteParam := c.Query("favorite")

	var userItems []models.UserItem
	var total int64
	var err error

	// Filtrar por favoritos
	if favoriteParam == "true" {
		userItems, total, err = h.service.GetMyFavorites(ctx, userID, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response := dto.NewPaginatedResponse(userItems, params.Page, params.Limit, total)
		c.JSON(http.StatusOK, response)
		return
	}

	// Filtrar por status
	if statusParam != "" {
		status := models.MediaStatus(statusParam)
		userItems, total, err = h.service.GetMyListByStatus(ctx, userID, status, params)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := dto.NewPaginatedResponse(userItems, params.Page, params.Limit, total)
		c.JSON(http.StatusOK, response)
		return
	}

	// Retornar lista completa
	userItems, total, err = h.service.GetMyList(ctx, userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.NewPaginatedResponse(userItems, params.Page, params.Limit, total)
	c.JSON(http.StatusOK, response)
}

// GetMyListItem retorna um item específico da lista do usuário
// @Summary      Get list item
// @Description  Get a specific item from user's list by ID
// @Tags         my-list
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "User Item ID"
// @Success      200  {object}  models.UserItem       "Success - returns user item"
// @Failure      400  {object}  map[string]string     "Bad request - invalid ID"
// @Failure      404  {object}  map[string]string     "Item not found"
// @Router       /my-list/{id} [get]
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
// @Summary      Update list item
// @Description  Update a user's list item (status, rating, progress, etc)
// @Tags         my-list
// @Accept       json
// @Produce      json
// @Param        id    path  int                  true  "User Item ID"
// @Param        item  body  models.UserItem      true  "Updated item data"
// @Success      200  {object}  models.UserItem       "Item updated successfully"
// @Failure      400  {object}  map[string]string     "Bad request - validation error"
// @Failure      404  {object}  map[string]string     "Item not found"
// @Router       /my-list/{id} [put]
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
// @Summary      Remove from list
// @Description  Remove an item from user's tracking list
// @Tags         my-list
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "User Item ID"
// @Success      204  "Item removed successfully"
// @Failure      400  {object}  map[string]string  "Bad request - invalid ID"
// @Failure      404  {object}  map[string]string  "Item not found"
// @Router       /my-list/{id} [delete]
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
// @Summary      Get list statistics
// @Description  Get statistics about user's tracking list (totals by status, favorites count)
// @Tags         my-list
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]int64       "Success - returns statistics"
// @Failure      500  {object}  map[string]string      "Internal server error"
// @Router       /my-list/stats [get]
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
