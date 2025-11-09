package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/services"
	"gorm.io/gorm"
)

type TagHandler struct {
	tagService *services.TagService
}

// NewTagHandler cria uma nova instância do handler de tags
func NewTagHandler(tagService *services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// CreateTagRequest representa o payload de criação de tag
type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateTagRequest representa o payload de atualização de tag
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateTag godoc
// @Summary Criar uma nova tag
// @Description Cria uma nova tag para categorizar items
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body CreateTagRequest true "Dados da tag"
// @Success 201 {object} models.Tag
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &models.Tag{
		Name: req.Name,
	}

	if err := h.tagService.CreateTag(ctx, tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GetAllTags godoc
// @Summary Listar todas as tags
// @Description Retorna uma lista de todas as tags
// @Tags tags
// @Produce json
// @Success 200 {array} models.Tag
// @Failure 500 {object} map[string]string
// @Router /api/tags [get]
func (h *TagHandler) GetAllTags(c *gin.Context) {
	ctx := c.Request.Context()

	tags, err := h.tagService.GetAllTags(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// GetTagByID godoc
// @Summary Buscar tag por ID
// @Description Retorna uma tag específica pelo seu ID
// @Tags tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} models.Tag
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/tags/{id} [get]
func (h *TagHandler) GetTagByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	tag, err := h.tagService.GetTagByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tag"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// UpdateTag godoc
// @Summary Atualizar uma tag
// @Description Atualiza os dados de uma tag existente
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param tag body UpdateTagRequest true "Dados atualizados"
// @Success 200 {object} models.Tag
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &models.Tag{
		Name: req.Name,
	}

	if err := h.tagService.UpdateTag(ctx, uint(id), tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar a tag atualizada
	updatedTag, err := h.tagService.GetTagByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated tag"})
		return
	}

	c.JSON(http.StatusOK, updatedTag)
}

// DeleteTag godoc
// @Summary Deletar uma tag
// @Description Remove uma tag do banco de dados
// @Tags tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := h.tagService.DeleteTag(ctx, uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}
