package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/dto"
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

// CreateTag cria uma nova tag
// @Summary      Create tag
// @Description  Create a new tag for categorizing items
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        tag  body  CreateTagRequest  true  "Tag data"
// @Success      201  {object}  models.Tag            "Tag created successfully"
// @Failure      400  {object}  map[string]string     "Bad request - validation error"
// @Failure      500  {object}  map[string]string     "Internal server error"
// @Router       /tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateTagRequest
	if err := validateAndBind(c, &req); err != nil {
		respondValidationError(c, err)
		return
	}

	tag := &models.Tag{
		Name: req.Name,
	}

	if err := h.tagService.CreateTag(ctx, tag); err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, tag)
}

// GetAllTags retorna todas as tags
// @Summary      Get all tags
// @Description  Get a list of all available tags
// @Tags         tags
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Tag            "Success - returns all tags"
// @Failure      500  {object}  map[string]string     "Internal server error"
// @Router       /tags [get]
func (h *TagHandler) GetAllTags(c *gin.Context) {
	ctx := c.Request.Context()

	tags, err := h.tagService.GetAllTags(ctx)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, tags)
}

// GetTagByID retorna uma tag específica
// @Summary      Get tag by ID
// @Description  Get a specific tag by its ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Tag ID"
// @Success      200  {object}  models.Tag            "Success - returns tag"
// @Failure      400  {object}  map[string]string     "Bad request - invalid ID"
// @Failure      404  {object}  map[string]string     "Tag not found"
// @Router       /tags/{id} [get]
func (h *TagHandler) GetTagByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := validateID(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeInvalidID, err.Error())
		return
	}

	tag, err := h.tagService.GetTagByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			respondNotFound(c, "Tag")
			return
		}
		respondInternalError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, tag)
}

// UpdateTag atualiza uma tag existente
// @Summary      Update tag
// @Description  Update an existing tag's data
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path  int                true  "Tag ID"
// @Param        tag  body  UpdateTagRequest   true  "Updated tag data"
// @Success      200  {object}  models.Tag            "Tag updated successfully"
// @Failure      400  {object}  map[string]string     "Bad request - validation error"
// @Failure      404  {object}  map[string]string     "Tag not found"
// @Router       /tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := validateID(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeInvalidID, err.Error())
		return
	}

	var req UpdateTagRequest
	if err := validateAndBind(c, &req); err != nil {
		respondValidationError(c, err)
		return
	}

	tag := &models.Tag{
		Name: req.Name,
	}

	if err := h.tagService.UpdateTag(ctx, id, tag); err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, err.Error())
		return
	}

	// Buscar a tag atualizada
	updatedTag, err := h.tagService.GetTagByID(ctx, id)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, updatedTag)
}

// DeleteTag remove uma tag
// @Summary      Delete tag
// @Description  Delete a tag from the database
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Tag ID"
// @Success      200  {object}  map[string]string  "Tag deleted successfully"
// @Failure      400  {object}  map[string]string  "Bad request - invalid ID"
// @Failure      404  {object}  map[string]string  "Tag not found"
// @Router       /tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := validateID(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeInvalidID, err.Error())
		return
	}

	if err := h.tagService.DeleteTag(ctx, id); err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}
