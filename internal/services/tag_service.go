package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/repositories"
	"gorm.io/gorm"
)

type TagService struct {
	tagRepo repositories.TagRepositoryInterface
}

// NewTagService cria uma nova instância do serviço de tags
func NewTagService(tagRepo repositories.TagRepositoryInterface) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// CreateTag cria uma nova tag com validações
func (s *TagService) CreateTag(ctx context.Context, tag *models.Tag) error {
	// Validar dados
	if err := s.validateTag(tag); err != nil {
		return err
	}

	// Normalizar nome (lowercase, trim)
	tag.Name = s.normalizeName(tag.Name)

	// Verificar se já existe
	existing, err := s.tagRepo.GetByName(ctx, tag.Name)
	if err == nil && existing != nil {
		return models.ErrDuplicateTag
	}

	// Criar tag
	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

// GetAllTags retorna todas as tags
func (s *TagService) GetAllTags(ctx context.Context) ([]models.Tag, error) {
	return s.tagRepo.GetAll(ctx)
}

// GetTagByID retorna uma tag específica
func (s *TagService) GetTagByID(ctx context.Context, id uint) (*models.Tag, error) {
	if id == 0 {
		return nil, errors.New("invalid tag ID")
	}
	return s.tagRepo.GetByID(ctx, id)
}

// GetTagByName retorna uma tag pelo nome
func (s *TagService) GetTagByName(ctx context.Context, name string) (*models.Tag, error) {
	if name == "" {
		return nil, errors.New("tag name is required")
	}
	name = s.normalizeName(name)
	return s.tagRepo.GetByName(ctx, name)
}

// UpdateTag atualiza uma tag existente
func (s *TagService) UpdateTag(ctx context.Context, id uint, updatedTag *models.Tag) error {
	// Verificar se a tag existe
	existingTag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("tag not found")
		}
		return fmt.Errorf("failed to find tag: %w", err)
	}

	// Validar dados atualizados
	if err := s.validateTag(updatedTag); err != nil {
		return err
	}

	// Normalizar nome
	updatedTag.Name = s.normalizeName(updatedTag.Name)

	// Verificar se o novo nome já existe em outra tag
	if updatedTag.Name != existingTag.Name {
		existing, err := s.tagRepo.GetByName(ctx, updatedTag.Name)
		if err == nil && existing != nil && existing.ID != id {
			return models.ErrDuplicateTag
		}
	}

	// Manter o ID original
	updatedTag.ID = id

	// Atualizar tag
	if err := s.tagRepo.Update(ctx, updatedTag); err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	return nil
}

// DeleteTag remove uma tag
func (s *TagService) DeleteTag(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid tag ID")
	}

	// Verificar se a tag existe
	if _, err := s.tagRepo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("tag not found")
		}
		return fmt.Errorf("failed to find tag: %w", err)
	}

	return s.tagRepo.Delete(ctx, id)
}

// FindOrCreateTags busca ou cria tags pelo nome
func (s *TagService) FindOrCreateTags(ctx context.Context, names []string) ([]models.Tag, error) {
	var tags []models.Tag
	for _, name := range names {
		name = s.normalizeName(name)
		if name == "" {
			continue
		}

		tag := &models.Tag{Name: name}
		err := s.tagRepo.FindOrCreate(ctx, tag)
		if err != nil {
			return nil, fmt.Errorf("failed to find or create tag '%s': %w", name, err)
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}

// validateTag valida os dados de uma tag
func (s *TagService) validateTag(tag *models.Tag) error {
	if tag.Name == "" {
		return errors.New("tag name is required")
	}

	if len(tag.Name) < 2 {
		return errors.New("tag name must have at least 2 characters")
	}

	if len(tag.Name) > 50 {
		return errors.New("tag name must have at most 50 characters")
	}

	return nil
}

// normalizeName normaliza o nome da tag (lowercase e trim)
func (s *TagService) normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
