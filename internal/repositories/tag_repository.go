package repositories

import (
	"context"

	"github.com/rafaelc-rb/geekery-api/internal/models"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository cria uma nova instância do repositório de tags
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create cria uma nova tag no banco de dados
func (r *TagRepository) Create(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// GetAll retorna todas as tags
func (r *TagRepository) GetAll(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.WithContext(ctx).Find(&tags).Error
	return tags, err
}

// GetByID retorna uma tag específica pelo ID
func (r *TagRepository) GetByID(ctx context.Context, id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.WithContext(ctx).Preload("Items").First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetByName retorna uma tag pelo nome
func (r *TagRepository) GetByName(ctx context.Context, name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// Update atualiza uma tag existente
func (r *TagRepository) Update(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

// Delete remove uma tag do banco de dados (soft delete)
func (r *TagRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Tag{}, id).Error
}

// FindOrCreate busca uma tag pelo nome ou cria se não existir
func (r *TagRepository) FindOrCreate(ctx context.Context, tag *models.Tag) error {
	err := r.db.WithContext(ctx).Where("name = ?", tag.Name).FirstOrCreate(tag).Error
	if err != nil {
		return err
	}
	return nil
}

// GetTagsByIDs retorna múltiplas tags pelos seus IDs
func (r *TagRepository) GetTagsByIDs(ctx context.Context, ids []uint) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.WithContext(ctx).Find(&tags, ids).Error
	return tags, err
}
