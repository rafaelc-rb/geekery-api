package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// GetAllItems retorna todos os items do catálogo
// GET /api/items
func (h *ItemHandler) GetAllItems(c *gin.Context) {
	ctx := c.Request.Context()

	// Parâmetro de filtro opcional por tipo
	typeParam := c.Query("type")

	var items []models.Item
	var err error

	if typeParam != "" {
		mediaType := models.MediaType(typeParam)
		items, err = h.service.GetItemsByType(ctx, mediaType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		items, err = h.service.GetAllItems(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, items)
}

// GetItemByID retorna um item específico do catálogo
// GET /api/items/:id
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

// SearchItems busca items por título
// GET /api/items/search?q=query
func (h *ItemHandler) SearchItems(c *gin.Context) {
	ctx := c.Request.Context()
	query := c.Query("q")

	items, err := h.service.SearchItems(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// CreateItem cria um novo item no catálogo (admin apenas - futuro)
// POST /api/items
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
// PUT /api/items/:id
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
// DELETE /api/items/:id
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
