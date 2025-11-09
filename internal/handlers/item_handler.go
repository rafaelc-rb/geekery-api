package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/dto"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/services"
)

type ItemHandler struct {
	service *services.ItemService
}

// NewItemHandler cria uma nova instância do handler de items (catálogo global)
func NewItemHandler(service *services.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

// GetAllItems retorna todos os items do catálogo com paginação
// @Summary      Get all items
// @Description  Get paginated list of all items from the global catalog
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        page   query  int     false  "Page number" default(1) minimum(1)
// @Param        limit  query  int     false  "Items per page" default(20) minimum(1) maximum(100)
// @Param        type   query  string  false  "Filter by media type" Enums(anime, movie, series, game, manga, light_novel, music, book)
// @Success      200  {object}  dto.PaginatedResponse  "Success - returns paginated items"
// @Failure      400  {object}  map[string]string      "Bad request - invalid parameters"
// @Failure      500  {object}  map[string]string      "Internal server error"
// @Router       /items [get]
func (h *ItemHandler) GetAllItems(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse parâmetros de paginação
	var params dto.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}
	params.Normalize()

	// Parâmetro de filtro opcional por tipo
	typeParam := c.Query("type")

	var items []models.Item
	var total int64
	var err error

	if typeParam != "" {
		mediaType := models.MediaType(typeParam)
		items, total, err = h.service.GetItemsByType(ctx, mediaType, params)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		items, total, err = h.service.GetAllItems(ctx, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Retornar resposta paginada
	response := dto.NewPaginatedResponse(items, params.Page, params.Limit, total)
	c.JSON(http.StatusOK, response)
}

// GetItemByID retorna um item específico do catálogo
// @Summary      Get item by ID
// @Description  Get a specific item from the catalog by its ID
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Item ID"
// @Success      200  {object}  models.Item           "Success - returns item"
// @Failure      400  {object}  map[string]string     "Bad request - invalid ID"
// @Failure      404  {object}  map[string]string     "Item not found"
// @Router       /items/{id} [get]
func (h *ItemHandler) GetItemByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	item, err := h.service.GetItemByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// SearchItems busca items por título com paginação
// @Summary      Search items
// @Description  Search items by title (case-insensitive) with pagination
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        q      query  string  false  "Search query"
// @Param        page   query  int     false  "Page number" default(1)
// @Param        limit  query  int     false  "Items per page" default(20)
// @Success      200  {object}  dto.PaginatedResponse  "Success - returns matching items"
// @Failure      500  {object}  map[string]string      "Internal server error"
// @Router       /items/search [get]
func (h *ItemHandler) SearchItems(c *gin.Context) {
	ctx := c.Request.Context()
	query := c.Query("q")

	// Parse parâmetros de paginação
	var params dto.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}
	params.Normalize()

	items, total, err := h.service.SearchItems(ctx, query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar resposta paginada
	response := dto.NewPaginatedResponse(items, params.Page, params.Limit, total)
	c.JSON(http.StatusOK, response)
}

// CreateItem cria um novo item no catálogo (admin apenas - futuro)
// @Summary      Create item
// @Description  Create a new item in the global catalog (admin only in the future)
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        item  body  models.Item  true  "Item to create"
// @Success      201  {object}  models.Item           "Item created successfully"
// @Failure      400  {object}  map[string]string     "Bad request - validation error"
// @Router       /items [post]
func (h *ItemHandler) CreateItem(c *gin.Context) {
	ctx := c.Request.Context()

	var input struct {
		models.Item
		TagIDs []uint `json:"tag_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateItem(ctx, &input.Item, input.TagIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input.Item)
}

// UpdateItem atualiza um item do catálogo (admin apenas - futuro)
// @Summary      Update item
// @Description  Update an existing item in the catalog (admin only in the future)
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id    path  int          true  "Item ID"
// @Param        item  body  models.Item  true  "Item data to update"
// @Success      200  {object}  map[string]string  "Item updated successfully"
// @Failure      400  {object}  map[string]string  "Bad request - validation error"
// @Failure      404  {object}  map[string]string  "Item not found"
// @Router       /items/{id} [put]
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	var input struct {
		models.Item
		TagIDs []uint `json:"tag_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateItem(ctx, uint(id), &input.Item, input.TagIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item updated successfully"})
}

// DeleteItem remove um item do catálogo (admin apenas - futuro)
// @Summary      Delete item
// @Description  Delete an item from the catalog (admin only in the future)
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Item ID"
// @Success      204  "Item deleted successfully"
// @Failure      400  {object}  map[string]string  "Bad request - invalid ID"
// @Failure      404  {object}  map[string]string  "Item not found"
// @Router       /items/{id} [delete]
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	if err := h.service.DeleteItem(ctx, uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
