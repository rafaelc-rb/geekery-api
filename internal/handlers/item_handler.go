package handlers

import (
	"net/http"
	"strings"

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
		respondValidationError(c, err)
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
			respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, err.Error())
			return
		}
	} else {
		items, total, err = h.service.GetAllItems(ctx, params)
		if err != nil {
			respondInternalError(c, err)
			return
		}
	}

	// Retornar resposta paginada
	response := dto.NewPaginatedResponse(items, params.Page, params.Limit, total)
	respondSuccess(c, http.StatusOK, response)
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

	id, err := validateID(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeInvalidID, err.Error())
		return
	}

	item, err := h.service.GetItemByID(ctx, id)
	if err != nil {
		respondNotFound(c, "Item")
		return
	}

	respondSuccess(c, http.StatusOK, item)
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
		respondValidationError(c, err)
		return
	}
	params.Normalize()

	items, total, err := h.service.SearchItems(ctx, query, params)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	// Retornar resposta paginada
	response := dto.NewPaginatedResponse(items, params.Page, params.Limit, total)
	respondSuccess(c, http.StatusOK, response)
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

	if err := validateAndBind(c, &input); err != nil {
		respondValidationError(c, err)
		return
	}

	if err := h.service.CreateItem(ctx, &input.Item, input.TagIDs); err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, input.Item)
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

	id, err := validateID(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeInvalidID, err.Error())
		return
	}

	var input struct {
		models.Item
		TagIDs []uint `json:"tag_ids"`
	}

	if err := validateAndBind(c, &input); err != nil {
		respondValidationError(c, err)
		return
	}

	if err := h.service.UpdateItem(ctx, id, &input.Item, input.TagIDs); err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"message": "item updated successfully"})
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

	id, err := validateID(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeInvalidID, err.Error())
		return
	}

	if err := h.service.DeleteItem(ctx, id); err != nil {
		respondNotFound(c, "Item")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ImportAnime importa múltiplos animes de um arquivo CSV
// @Summary      Import anime items
// @Description  Import multiple anime items from CSV file
// @Tags         items
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV file with anime data"
// @Success      200  {object}  dto.ImportResult  "Import completed"
// @Failure      400  {object}  map[string]string  "Bad request"
// @Router       /items/import/anime [post]
func (h *ItemHandler) ImportAnime(c *gin.Context) {
	h.importByType(c, models.MediaTypeAnime)
}

// ImportManga importa múltiplos mangas de um arquivo CSV
// @Summary      Import manga items
// @Description  Import multiple manga items from CSV file
// @Tags         items
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV file with manga data"
// @Success      200  {object}  dto.ImportResult  "Import completed"
// @Failure      400  {object}  map[string]string  "Bad request"
// @Router       /items/import/manga [post]
func (h *ItemHandler) ImportManga(c *gin.Context) {
	h.importByType(c, models.MediaTypeManga)
}

// ImportMovie importa múltiplos filmes de um arquivo CSV
// @Summary      Import movie items
// @Description  Import multiple movie items from CSV file
// @Tags         items
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV file with movie data"
// @Success      200  {object}  dto.ImportResult  "Import completed"
// @Failure      400  {object}  map[string]string  "Bad request"
// @Router       /items/import/movie [post]
func (h *ItemHandler) ImportMovie(c *gin.Context) {
	h.importByType(c, models.MediaTypeMovie)
}

// ImportSeries importa múltiplas séries de um arquivo CSV
// @Summary      Import series items
// @Description  Import multiple series items from CSV file
// @Tags         items
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV file with series data"
// @Success      200  {object}  dto.ImportResult  "Import completed"
// @Failure      400  {object}  map[string]string  "Bad request"
// @Router       /items/import/series [post]
func (h *ItemHandler) ImportSeries(c *gin.Context) {
	h.importByType(c, models.MediaTypeSeries)
}

// ImportGame importa múltiplos games de um arquivo CSV
// @Summary      Import game items
// @Description  Import multiple game items from CSV file
// @Tags         items
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV file with game data"
// @Success      200  {object}  dto.ImportResult  "Import completed"
// @Failure      400  {object}  map[string]string  "Bad request"
// @Router       /items/import/game [post]
func (h *ItemHandler) ImportGame(c *gin.Context) {
	h.importByType(c, models.MediaTypeGame)
}

// ImportBook importa múltiplos livros de um arquivo CSV
// @Summary      Import book items
// @Description  Import multiple book items from CSV file
// @Tags         items
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV file with book data"
// @Success      200  {object}  dto.ImportResult  "Import completed"
// @Failure      400  {object}  map[string]string  "Bad request"
// @Router       /items/import/book [post]
func (h *ItemHandler) ImportBook(c *gin.Context) {
	h.importByType(c, models.MediaTypeBook)
}

// ImportMusic importa múltiplas músicas de um arquivo CSV
// @Summary      Import music items
// @Description  Import multiple music items from CSV file
// @Tags         items
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV file with music data"
// @Success      200  {object}  dto.ImportResult  "Import completed"
// @Failure      400  {object}  map[string]string  "Bad request"
// @Router       /items/import/music [post]
func (h *ItemHandler) ImportMusic(c *gin.Context) {
	h.importByType(c, models.MediaTypeMusic)
}

// importByType é o handler genérico de import por tipo
func (h *ItemHandler) importByType(c *gin.Context, mediaType models.MediaType) {
	ctx := c.Request.Context()

	// Receber arquivo
	file, err := c.FormFile("file")
	if err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, "CSV file is required")
		return
	}

	// Validar extensão
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, "File must be a .csv file")
		return
	}

	// Validar tamanho (max 10MB)
	if file.Size > 10*1024*1024 {
		respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, "File size must be less than 10MB")
		return
	}

	// Abrir arquivo
	src, err := file.Open()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	defer src.Close()

	// Processar no service
	result, err := h.service.ImportItemsFromCSV(ctx, src, mediaType)
	if err != nil {
		respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, result)
}
