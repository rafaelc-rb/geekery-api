package repositories

import (
	"context"
	"errors"

	"github.com/rafaelc-rb/geekery-api/internal/models"
	"gorm.io/gorm"
)

type UserItemRepository struct {
	db *gorm.DB
}

// NewUserItemRepository cria uma nova instância do repositório de user items
func NewUserItemRepository(db *gorm.DB) *UserItemRepository {
	return &UserItemRepository{db: db}
}

// Create adiciona um item à lista do usuário
func (r *UserItemRepository) Create(ctx context.Context, userItem *models.UserItem) error {
	return r.db.WithContext(ctx).Create(userItem).Error
}

// GetByUserID retorna todos os items da lista de um usuário
func (r *UserItemRepository) GetByUserID(ctx context.Context, userID uint) ([]models.UserItem, error) {
	var userItems []models.UserItem
	err := r.db.WithContext(ctx).Preload("Item").Preload("Item.Tags").Where("user_id = ?", userID).Find(&userItems).Error
	return userItems, err
}

// GetByUserAndItem busca um item específico na lista do usuário
func (r *UserItemRepository) GetByUserAndItem(ctx context.Context, userID, itemID uint) (*models.UserItem, error) {
	var userItem models.UserItem
	err := r.db.WithContext(ctx).Preload("Item").Preload("Item.Tags").Where("user_id = ? AND item_id = ?", userID, itemID).First(&userItem).Error
	if err != nil {
		return nil, err
	}
	return &userItem, nil
}

// GetByID retorna um user item pelo ID
func (r *UserItemRepository) GetByID(ctx context.Context, id uint) (*models.UserItem, error) {
	var userItem models.UserItem
	err := r.db.WithContext(ctx).Preload("Item").Preload("Item.Tags").First(&userItem, id).Error
	if err != nil {
		return nil, err
	}
	return &userItem, nil
}

// Update atualiza um item da lista do usuário
func (r *UserItemRepository) Update(ctx context.Context, userItem *models.UserItem) error {
	return r.db.WithContext(ctx).Save(userItem).Error
}

// Delete remove um item da lista do usuário
func (r *UserItemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.UserItem{}, id).Error
}

// Exists verifica se um item já está na lista do usuário
func (r *UserItemRepository) Exists(ctx context.Context, userID, itemID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.UserItem{}).Where("user_id = ? AND item_id = ?", userID, itemID).Count(&count).Error
	return count > 0, err
}

// GetByStatus retorna items do usuário filtrados por status
func (r *UserItemRepository) GetByStatus(ctx context.Context, userID uint, status models.MediaStatus) ([]models.UserItem, error) {
	var userItems []models.UserItem
	err := r.db.WithContext(ctx).Preload("Item").Preload("Item.Tags").
		Where("user_id = ? AND status = ?", userID, status).
		Find(&userItems).Error
	return userItems, err
}

// GetFavorites retorna todos os items favoritos do usuário
func (r *UserItemRepository) GetFavorites(ctx context.Context, userID uint) ([]models.UserItem, error) {
	var userItems []models.UserItem
	err := r.db.WithContext(ctx).Preload("Item").Preload("Item.Tags").
		Where("user_id = ? AND favorite = ?", userID, true).
		Find(&userItems).Error
	return userItems, err
}

// GetStatistics retorna estatísticas da lista do usuário
func (r *UserItemRepository) GetStatistics(ctx context.Context, userID uint) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Total de items
	var total int64
	err := r.db.WithContext(ctx).Model(&models.UserItem{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	// Por status
	statuses := []models.MediaStatus{
		models.StatusInProgress,
		models.StatusCompleted,
		models.StatusPlanned,
		models.StatusPaused,
		models.StatusDropped,
	}

	for _, status := range statuses {
		var count int64
		err := r.db.WithContext(ctx).Model(&models.UserItem{}).
			Where("user_id = ? AND status = ?", userID, status).
			Count(&count).Error
		if err != nil {
			return nil, err
		}
		stats[string(status)] = count
	}

	// Favoritos
	var favorites int64
	err = r.db.WithContext(ctx).Model(&models.UserItem{}).
		Where("user_id = ? AND favorite = ?", userID, true).
		Count(&favorites).Error
	if err != nil {
		return nil, err
	}
	stats["favorites"] = favorites

	return stats, nil
}

// GetByIDAndUser busca um user item por ID garantindo que pertence ao usuário
func (r *UserItemRepository) GetByIDAndUser(ctx context.Context, id, userID uint) (*models.UserItem, error) {
	var userItem models.UserItem
	err := r.db.WithContext(ctx).Preload("Item").Preload("Item.Tags").
		Where("id = ? AND user_id = ?", id, userID).
		First(&userItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user item not found or doesn't belong to user")
		}
		return nil, err
	}
	return &userItem, nil
}
