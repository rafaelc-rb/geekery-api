package repositories

import (
	"context"

	"github.com/rafaelc-rb/geekery-api/internal/dto"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

// NewItemRepository cria uma nova instância do repositório de items
func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// Create cria um novo item no catálogo global
func (r *ItemRepository) Create(ctx context.Context, item *models.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// GetAll retorna todos os items do catálogo com paginação
func (r *ItemRepository) GetAll(ctx context.Context, params dto.PaginationParams) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	// Normalizar parâmetros
	params.Normalize()

	// Contar total de items
	if err := r.db.WithContext(ctx).Model(&models.Item{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar items paginados
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Limit(params.Limit).
		Offset(params.GetOffset()).
		Find(&items).Error

	return items, total, err
}

// GetByID retorna um item específico pelo ID com Preload condicional baseado no tipo
func (r *ItemRepository) GetByID(ctx context.Context, id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.WithContext(ctx).Preload("Tags").First(&item, id).Error
	if err != nil {
		return nil, err
	}

	// Preload dados específicos baseado no tipo
	if err := r.preloadSpecificData(ctx, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// preloadSpecificData carrega os dados específicos baseado no tipo do item
func (r *ItemRepository) preloadSpecificData(ctx context.Context, item *models.Item) error {
	switch item.Type {
	case models.MediaTypeAnime:
		var animeData models.AnimeData
		err := r.db.WithContext(ctx).Where("item_id = ?", item.ID).First(&animeData).Error
		if err == nil {
			item.AnimeData = &animeData
		}
	case models.MediaTypeMovie:
		var movieData models.MovieData
		err := r.db.WithContext(ctx).Where("item_id = ?", item.ID).First(&movieData).Error
		if err == nil {
			item.MovieData = &movieData
		}
	case models.MediaTypeGame:
		var gameData models.GameData
		err := r.db.WithContext(ctx).Where("item_id = ?", item.ID).First(&gameData).Error
		if err == nil {
			item.GameData = &gameData
		}
	case models.MediaTypeMusic:
		var musicData models.MusicData
		err := r.db.WithContext(ctx).Where("item_id = ?", item.ID).First(&musicData).Error
		if err == nil {
			item.MusicData = &musicData
		}
	case models.MediaTypeBook:
		var bookData models.BookData
		err := r.db.WithContext(ctx).Where("item_id = ?", item.ID).First(&bookData).Error
		if err == nil {
			item.BookData = &bookData
		}
	case models.MediaTypeSeries:
		var seriesData models.SeriesData
		err := r.db.WithContext(ctx).Where("item_id = ?", item.ID).First(&seriesData).Error
		if err == nil {
			item.SeriesData = &seriesData
		}
	case models.MediaTypeManga, models.MediaTypeLightNovel:
		var bookData models.BookData
		err := r.db.WithContext(ctx).Where("item_id = ?", item.ID).First(&bookData).Error
		if err == nil {
			item.BookData = &bookData
		}
	}

	return nil
}

// GetByType retorna items filtrados por tipo com dados específicos e paginação
func (r *ItemRepository) GetByType(ctx context.Context, mediaType models.MediaType, params dto.PaginationParams) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	// Normalizar parâmetros
	params.Normalize()

	// Contar total de items do tipo
	if err := r.db.WithContext(ctx).Model(&models.Item{}).Where("type = ?", mediaType).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar items paginados
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Where("type = ?", mediaType).
		Limit(params.Limit).
		Offset(params.GetOffset()).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	// Preload dados específicos para todos os items do tipo
	for i := range items {
		if err := r.preloadSpecificData(ctx, &items[i]); err != nil {
			return nil, 0, err
		}
	}

	return items, total, nil
}

// Update atualiza um item existente no catálogo
func (r *ItemRepository) Update(ctx context.Context, item *models.Item) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// Delete remove um item do catálogo
func (r *ItemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Item{}, id).Error
}

// SearchByTitle busca items por título (case-insensitive) com paginação
func (r *ItemRepository) SearchByTitle(ctx context.Context, query string, params dto.PaginationParams) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	// Normalizar parâmetros
	params.Normalize()

	// Contar total de resultados
	searchQuery := "%" + query + "%"
	if err := r.db.WithContext(ctx).Model(&models.Item{}).
		Where("LOWER(title) LIKE LOWER(?)", searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar items paginados
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Where("LOWER(title) LIKE LOWER(?)", searchQuery).
		Limit(params.Limit).
		Offset(params.GetOffset()).
		Find(&items).Error

	return items, total, err
}

// GetByExternalID busca um item por ID externo (MAL, IMDb, etc)
func (r *ItemRepository) GetByExternalID(ctx context.Context, source, externalID string) (*models.Item, error) {
	var item models.Item
	// Busca no campo JSONB external_metadata
	err := r.db.WithContext(ctx).Preload("Tags").
		Where("external_metadata->>? = ?", source, externalID).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetByYear retorna items de um ano específico
func (r *ItemRepository) GetByYear(ctx context.Context, year int) ([]models.Item, error) {
	var items []models.Item
	// Busca items onde o ano da release_date é igual ao ano fornecido
	err := r.db.WithContext(ctx).Preload("Tags").
		Where("EXTRACT(YEAR FROM release_date) = ?", year).
		Find(&items).Error
	return items, err
}

// AssociateTags associa tags a um item
func (r *ItemRepository) AssociateTags(ctx context.Context, itemID uint, tagIDs []uint) error {
	var item models.Item
	if err := r.db.WithContext(ctx).First(&item, itemID).Error; err != nil {
		return err
	}

	var tags []models.Tag
	if err := r.db.WithContext(ctx).Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&item).Association("Tags").Replace(tags)
}

// RemoveTag remove uma tag específica de um item
func (r *ItemRepository) RemoveTag(ctx context.Context, itemID uint, tagID uint) error {
	var item models.Item
	if err := r.db.WithContext(ctx).First(&item, itemID).Error; err != nil {
		return err
	}

	var tag models.Tag
	if err := r.db.WithContext(ctx).First(&tag, tagID).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&item).Association("Tags").Delete(&tag)
}

// CreateSpecificData cria dados específicos para um item baseado no tipo
func (r *ItemRepository) CreateSpecificData(ctx context.Context, itemID uint, mediaType models.MediaType, data interface{}) error {
	switch mediaType {
	case models.MediaTypeAnime:
		if animeData, ok := data.(*models.AnimeData); ok {
			animeData.ItemID = itemID
			return r.db.WithContext(ctx).Create(animeData).Error
		}
	case models.MediaTypeMovie:
		if movieData, ok := data.(*models.MovieData); ok {
			movieData.ItemID = itemID
			return r.db.WithContext(ctx).Create(movieData).Error
		}
	case models.MediaTypeGame:
		if gameData, ok := data.(*models.GameData); ok {
			gameData.ItemID = itemID
			return r.db.WithContext(ctx).Create(gameData).Error
		}
	case models.MediaTypeMusic:
		if musicData, ok := data.(*models.MusicData); ok {
			musicData.ItemID = itemID
			return r.db.WithContext(ctx).Create(musicData).Error
		}
	case models.MediaTypeBook:
		if bookData, ok := data.(*models.BookData); ok {
			bookData.ItemID = itemID
			return r.db.WithContext(ctx).Create(bookData).Error
		}
	case models.MediaTypeSeries:
		if seriesData, ok := data.(*models.SeriesData); ok {
			seriesData.ItemID = itemID
			return r.db.WithContext(ctx).Create(seriesData).Error
		}
	}
	return nil
}
