package repositories

import (
	"context"

	"github.com/rafaelc-rb/geekery-api/internal/dto"
	"github.com/rafaelc-rb/geekery-api/internal/models"
)

// ItemRepositoryInterface define os métodos do repositório de items
type ItemRepositoryInterface interface {
	Create(ctx context.Context, item *models.Item) error
	GetAll(ctx context.Context, params dto.PaginationParams) ([]models.Item, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Item, error)
	GetByType(ctx context.Context, mediaType models.MediaType, params dto.PaginationParams) ([]models.Item, int64, error)
	Update(ctx context.Context, item *models.Item) error
	Delete(ctx context.Context, id uint) error
	SearchByTitle(ctx context.Context, query string, params dto.PaginationParams) ([]models.Item, int64, error)
	GetByExternalID(ctx context.Context, source, externalID string) (*models.Item, error)
	GetByYear(ctx context.Context, year int) ([]models.Item, error)
	AssociateTags(ctx context.Context, itemID uint, tagIDs []uint) error
	RemoveTag(ctx context.Context, itemID uint, tagID uint) error
	CreateSpecificData(ctx context.Context, itemID uint, mediaType models.MediaType, data interface{}) error
}

// TagRepositoryInterface define os métodos do repositório de tags
type TagRepositoryInterface interface {
	Create(ctx context.Context, tag *models.Tag) error
	GetAll(ctx context.Context) ([]models.Tag, error)
	GetByID(ctx context.Context, id uint) (*models.Tag, error)
	GetByName(ctx context.Context, name string) (*models.Tag, error)
	Update(ctx context.Context, tag *models.Tag) error
	Delete(ctx context.Context, id uint) error
	FindOrCreate(ctx context.Context, tag *models.Tag) error
	GetTagsByIDs(ctx context.Context, ids []uint) ([]models.Tag, error)
}

// UserItemRepositoryInterface define os métodos do repositório de user_items
type UserItemRepositoryInterface interface {
	Create(ctx context.Context, userItem *models.UserItem) error
	GetByUserID(ctx context.Context, userID uint, params dto.PaginationParams) ([]models.UserItem, int64, error)
	GetByUserAndItem(ctx context.Context, userID uint, itemID uint) (*models.UserItem, error)
	GetByID(ctx context.Context, id uint) (*models.UserItem, error)
	Update(ctx context.Context, userItem *models.UserItem) error
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, userID uint, itemID uint) (bool, error)
	GetByStatus(ctx context.Context, userID uint, status models.MediaStatus, params dto.PaginationParams) ([]models.UserItem, int64, error)
	GetFavorites(ctx context.Context, userID uint, params dto.PaginationParams) ([]models.UserItem, int64, error)
	GetStatistics(ctx context.Context, userID uint) (map[string]int64, error)
	GetByIDAndUser(ctx context.Context, id uint, userID uint) (*models.UserItem, error)
}
