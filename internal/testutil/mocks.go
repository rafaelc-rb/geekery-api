package testutil

import (
	"context"

	"github.com/rafaelc-rb/geekery-api/internal/models"
)

// MockItemRepository é um mock do ItemRepository para testes
type MockItemRepository struct {
	CreateFunc              func(ctx context.Context, item *models.Item) error
	GetAllFunc              func(ctx context.Context) ([]models.Item, error)
	GetByIDFunc             func(ctx context.Context, id uint) (*models.Item, error)
	GetByTypeFunc           func(ctx context.Context, mediaType models.MediaType) ([]models.Item, error)
	UpdateFunc              func(ctx context.Context, item *models.Item) error
	DeleteFunc              func(ctx context.Context, id uint) error
	SearchByTitleFunc       func(ctx context.Context, query string) ([]models.Item, error)
	GetByExternalIDFunc     func(ctx context.Context, source, externalID string) (*models.Item, error)
	GetByYearFunc           func(ctx context.Context, year int) ([]models.Item, error)
	AssociateTagsFunc       func(ctx context.Context, itemID uint, tagIDs []uint) error
	RemoveTagFunc           func(ctx context.Context, itemID uint, tagID uint) error
	CreateSpecificDataFunc  func(ctx context.Context, itemID uint, mediaType models.MediaType, data interface{}) error
}

func (m *MockItemRepository) Create(ctx context.Context, item *models.Item) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, item)
	}
	return nil
}

func (m *MockItemRepository) GetAll(ctx context.Context) ([]models.Item, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return []models.Item{}, nil
}

func (m *MockItemRepository) GetByID(ctx context.Context, id uint) (*models.Item, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return &models.Item{}, nil
}

func (m *MockItemRepository) GetByType(ctx context.Context, mediaType models.MediaType) ([]models.Item, error) {
	if m.GetByTypeFunc != nil {
		return m.GetByTypeFunc(ctx, mediaType)
	}
	return []models.Item{}, nil
}

func (m *MockItemRepository) Update(ctx context.Context, item *models.Item) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, item)
	}
	return nil
}

func (m *MockItemRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockItemRepository) SearchByTitle(ctx context.Context, query string) ([]models.Item, error) {
	if m.SearchByTitleFunc != nil {
		return m.SearchByTitleFunc(ctx, query)
	}
	return []models.Item{}, nil
}

func (m *MockItemRepository) GetByExternalID(ctx context.Context, source, externalID string) (*models.Item, error) {
	if m.GetByExternalIDFunc != nil {
		return m.GetByExternalIDFunc(ctx, source, externalID)
	}
	return nil, nil
}

func (m *MockItemRepository) GetByYear(ctx context.Context, year int) ([]models.Item, error) {
	if m.GetByYearFunc != nil {
		return m.GetByYearFunc(ctx, year)
	}
	return []models.Item{}, nil
}

func (m *MockItemRepository) AssociateTags(ctx context.Context, itemID uint, tagIDs []uint) error {
	if m.AssociateTagsFunc != nil {
		return m.AssociateTagsFunc(ctx, itemID, tagIDs)
	}
	return nil
}

func (m *MockItemRepository) RemoveTag(ctx context.Context, itemID uint, tagID uint) error {
	if m.RemoveTagFunc != nil {
		return m.RemoveTagFunc(ctx, itemID, tagID)
	}
	return nil
}

func (m *MockItemRepository) CreateSpecificData(ctx context.Context, itemID uint, mediaType models.MediaType, data interface{}) error {
	if m.CreateSpecificDataFunc != nil {
		return m.CreateSpecificDataFunc(ctx, itemID, mediaType, data)
	}
	return nil
}

// MockUserItemRepository é um mock do UserItemRepository para testes
type MockUserItemRepository struct {
	CreateFunc          func(ctx context.Context, userItem *models.UserItem) error
	GetByUserIDFunc     func(ctx context.Context, userID uint) ([]models.UserItem, error)
	GetByIDFunc         func(ctx context.Context, id uint) (*models.UserItem, error)
	GetByIDAndUserFunc  func(ctx context.Context, id, userID uint) (*models.UserItem, error)
	UpdateFunc          func(ctx context.Context, userItem *models.UserItem) error
	DeleteFunc          func(ctx context.Context, id uint) error
	ExistsFunc          func(ctx context.Context, userID, itemID uint) (bool, error)
	GetByStatusFunc     func(ctx context.Context, userID uint, status models.MediaStatus) ([]models.UserItem, error)
	GetFavoritesFunc    func(ctx context.Context, userID uint) ([]models.UserItem, error)
	GetStatisticsFunc   func(ctx context.Context, userID uint) (map[string]int64, error)
}

func (m *MockUserItemRepository) Create(ctx context.Context, userItem *models.UserItem) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, userItem)
	}
	return nil
}

func (m *MockUserItemRepository) GetByUserID(ctx context.Context, userID uint) ([]models.UserItem, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(ctx, userID)
	}
	return []models.UserItem{}, nil
}

func (m *MockUserItemRepository) GetByID(ctx context.Context, id uint) (*models.UserItem, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return &models.UserItem{}, nil
}

func (m *MockUserItemRepository) GetByIDAndUser(ctx context.Context, id, userID uint) (*models.UserItem, error) {
	if m.GetByIDAndUserFunc != nil {
		return m.GetByIDAndUserFunc(ctx, id, userID)
	}
	return &models.UserItem{}, nil
}

func (m *MockUserItemRepository) Update(ctx context.Context, userItem *models.UserItem) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, userItem)
	}
	return nil
}

func (m *MockUserItemRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockUserItemRepository) Exists(ctx context.Context, userID, itemID uint) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, userID, itemID)
	}
	return false, nil
}

func (m *MockUserItemRepository) GetByStatus(ctx context.Context, userID uint, status models.MediaStatus) ([]models.UserItem, error) {
	if m.GetByStatusFunc != nil {
		return m.GetByStatusFunc(ctx, userID, status)
	}
	return []models.UserItem{}, nil
}

func (m *MockUserItemRepository) GetFavorites(ctx context.Context, userID uint) ([]models.UserItem, error) {
	if m.GetFavoritesFunc != nil {
		return m.GetFavoritesFunc(ctx, userID)
	}
	return []models.UserItem{}, nil
}

func (m *MockUserItemRepository) GetStatistics(ctx context.Context, userID uint) (map[string]int64, error) {
	if m.GetStatisticsFunc != nil {
		return m.GetStatisticsFunc(ctx, userID)
	}
	return map[string]int64{}, nil
}

func (m *MockUserItemRepository) GetByUserAndItem(ctx context.Context, userID uint, itemID uint) (*models.UserItem, error) {
	// Redirecionar para GetByIDAndUserFunc se existir
	if m.GetByIDAndUserFunc != nil {
		return m.GetByIDAndUserFunc(ctx, 0, userID) // ID não usado neste caso
	}
	return &models.UserItem{}, nil
}

// MockTagRepository é um mock do TagRepository para testes
type MockTagRepository struct {
	CreateFunc        func(ctx context.Context, tag *models.Tag) error
	GetAllFunc        func(ctx context.Context) ([]models.Tag, error)
	GetByIDFunc       func(ctx context.Context, id uint) (*models.Tag, error)
	GetByNameFunc     func(ctx context.Context, name string) (*models.Tag, error)
	UpdateFunc        func(ctx context.Context, tag *models.Tag) error
	DeleteFunc        func(ctx context.Context, id uint) error
	FindOrCreateFunc  func(ctx context.Context, tag *models.Tag) error
	GetTagsByIDsFunc  func(ctx context.Context, ids []uint) ([]models.Tag, error)
}

func (m *MockTagRepository) Create(ctx context.Context, tag *models.Tag) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, tag)
	}
	return nil
}

func (m *MockTagRepository) GetAll(ctx context.Context) ([]models.Tag, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return []models.Tag{}, nil
}

func (m *MockTagRepository) GetByID(ctx context.Context, id uint) (*models.Tag, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return &models.Tag{}, nil
}

func (m *MockTagRepository) GetByName(ctx context.Context, name string) (*models.Tag, error) {
	if m.GetByNameFunc != nil {
		return m.GetByNameFunc(ctx, name)
	}
	return &models.Tag{}, nil
}

func (m *MockTagRepository) Update(ctx context.Context, tag *models.Tag) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, tag)
	}
	return nil
}

func (m *MockTagRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockTagRepository) FindOrCreate(ctx context.Context, tag *models.Tag) error {
	if m.FindOrCreateFunc != nil {
		return m.FindOrCreateFunc(ctx, tag)
	}
	return nil
}

func (m *MockTagRepository) GetTagsByIDs(ctx context.Context, ids []uint) ([]models.Tag, error) {
	if m.GetTagsByIDsFunc != nil {
		return m.GetTagsByIDsFunc(ctx, ids)
	}
	return []models.Tag{}, nil
}
