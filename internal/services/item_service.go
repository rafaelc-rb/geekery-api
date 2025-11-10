package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rafaelc-rb/geekery-api/internal/dto"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/repositories"
	"gorm.io/gorm"
)

type ItemService struct {
	itemRepo repositories.ItemRepositoryInterface
	tagRepo  repositories.TagRepositoryInterface
}

// NewItemService cria uma nova instância do serviço de items (catálogo global)
func NewItemService(itemRepo repositories.ItemRepositoryInterface, tagRepo repositories.TagRepositoryInterface) *ItemService {
	return &ItemService{
		itemRepo: itemRepo,
		tagRepo:  tagRepo,
	}
}

// CreateItem cria um novo item no catálogo global (admin apenas)
func (s *ItemService) CreateItem(ctx context.Context, item *models.Item, tagIDs []uint, tagNames []string) error {
	// Validar
	if err := item.Validate(); err != nil {
		return err
	}

	// Criar o item base primeiro
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	// Processar tags por nome (find or create)
	if len(tagNames) > 0 {
		createdTagIDs, err := s.findOrCreateTags(ctx, tagNames)
		if err != nil {
			return err
		}
		tagIDs = append(tagIDs, createdTagIDs...)
	}

	// Associar tags se fornecidas
	if len(tagIDs) > 0 {
		if err := s.itemRepo.AssociateTags(ctx, item.ID, tagIDs); err != nil {
			return fmt.Errorf("failed to associate tags: %w", err)
		}
	}

	return nil
}

// CreateItemWithSpecificData cria um item com dados específicos baseado no tipo
func (s *ItemService) CreateItemWithSpecificData(ctx context.Context, item *models.Item, specificData interface{}, tagIDs []uint, tagNames []string) error {
	// Validar item base
	if err := item.Validate(); err != nil {
		return err
	}

	// Criar item base
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	// Criar dados específicos baseado no tipo
	if specificData != nil {
		if err := s.itemRepo.CreateSpecificData(ctx, item.ID, item.Type, specificData); err != nil {
			return fmt.Errorf("failed to create specific data: %w", err)
		}
	}

	// Processar tags por nome (find or create)
	if len(tagNames) > 0 {
		createdTagIDs, err := s.findOrCreateTags(ctx, tagNames)
		if err != nil {
			return err
		}
		tagIDs = append(tagIDs, createdTagIDs...)
	}

	// Associar tags
	if len(tagIDs) > 0 {
		if err := s.itemRepo.AssociateTags(ctx, item.ID, tagIDs); err != nil {
			return fmt.Errorf("failed to associate tags: %w", err)
		}
	}

	return nil
}

// GetAllItems retorna todos os items do catálogo com paginação
func (s *ItemService) GetAllItems(ctx context.Context, params dto.PaginationParams) ([]models.Item, int64, error) {
	return s.itemRepo.GetAll(ctx, params)
}

// GetItemByID retorna um item específico do catálogo
func (s *ItemService) GetItemByID(ctx context.Context, id uint) (*models.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	return item, nil
}

// GetItemsByType retorna items filtrados por tipo com paginação
func (s *ItemService) GetItemsByType(ctx context.Context, mediaType models.MediaType, params dto.PaginationParams) ([]models.Item, int64, error) {
	if !mediaType.IsValid() {
		return nil, 0, models.ErrInvalidMediaType
	}
	return s.itemRepo.GetByType(ctx, mediaType, params)
}

// SearchItems busca items por título com paginação
func (s *ItemService) SearchItems(ctx context.Context, query string, params dto.PaginationParams) ([]models.Item, int64, error) {
	if query == "" {
		return s.itemRepo.GetAll(ctx, params)
	}
	return s.itemRepo.SearchByTitle(ctx, query, params)
}

// UpdateItem atualiza um item do catálogo (admin apenas)
func (s *ItemService) UpdateItem(ctx context.Context, id uint, updatedItem *models.Item, tagIDs []uint, tagNames []string) error {
	// Verificar se existe
	existingItem, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item not found")
		}
		return fmt.Errorf("failed to get item: %w", err)
	}

	// Atualizar campos
	existingItem.Title = updatedItem.Title
	existingItem.Type = updatedItem.Type
	existingItem.Description = updatedItem.Description
	existingItem.ReleaseDate = updatedItem.ReleaseDate
	existingItem.CoverURL = updatedItem.CoverURL
	existingItem.ExternalMetadata = updatedItem.ExternalMetadata

	// Validar
	if err := existingItem.Validate(); err != nil {
		return err
	}

	// Salvar
	if err := s.itemRepo.Update(ctx, existingItem); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	// Processar tags por nome (find or create)
	if len(tagNames) > 0 {
		createdTagIDs, err := s.findOrCreateTags(ctx, tagNames)
		if err != nil {
			return err
		}
		tagIDs = append(tagIDs, createdTagIDs...)
	}

	// Atualizar tags se fornecidas
	if len(tagIDs) > 0 {
		if err := s.itemRepo.AssociateTags(ctx, id, tagIDs); err != nil {
			return fmt.Errorf("failed to update tags: %w", err)
		}
	}

	return nil
}

// DeleteItem remove um item do catálogo (admin apenas)
func (s *ItemService) DeleteItem(ctx context.Context, id uint) error {
	if err := s.itemRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item not found")
		}
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

// findOrCreateTags busca ou cria tags pelo nome e retorna seus IDs
func (s *ItemService) findOrCreateTags(ctx context.Context, tagNames []string) ([]uint, error) {
	tagIDs := []uint{}
	for _, name := range tagNames {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}

		tag := &models.Tag{Name: name}
		if err := s.tagRepo.FindOrCreate(ctx, tag); err != nil {
			return nil, fmt.Errorf("failed to find/create tag '%s': %w", name, err)
		}
		tagIDs = append(tagIDs, tag.ID)
	}
	return tagIDs, nil
}

// AssociateTags associa tags a um item do catálogo
func (s *ItemService) AssociateTags(ctx context.Context, itemID uint, tagIDs []uint) error {
	// Verificar se o item existe
	_, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item not found")
		}
		return fmt.Errorf("failed to verify item: %w", err)
	}

	// Associar tags
	if err := s.itemRepo.AssociateTags(ctx, itemID, tagIDs); err != nil {
		return fmt.Errorf("failed to associate tags: %w", err)
	}

	return nil
}
