package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/rafaelc-rb/geekery-api/internal/dto"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/repositories"
	"gorm.io/gorm"
)

type UserItemService struct {
	userItemRepo repositories.UserItemRepositoryInterface
	itemRepo     repositories.ItemRepositoryInterface
}

// NewUserItemService cria uma nova instância do serviço de user items
func NewUserItemService(userItemRepo repositories.UserItemRepositoryInterface, itemRepo repositories.ItemRepositoryInterface) *UserItemService {
	return &UserItemService{
		userItemRepo: userItemRepo,
		itemRepo:     itemRepo,
	}
}

// AddToList adiciona um item à lista do usuário
func (s *UserItemService) AddToList(ctx context.Context, userID uint, itemID uint, status models.MediaStatus) (*models.UserItem, error) {
	// Verificar se o item existe no catálogo
	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found in catalog")
		}
		return nil, fmt.Errorf("failed to verify item existence: %w", err)
	}

	// Verificar se o item já está na lista do usuário
	exists, err := s.userItemRepo.Exists(ctx, userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if item exists in user list: %w", err)
	}
	if exists {
		return nil, models.ErrDuplicateEntry
	}

	// Criar o user item com progress type padrão baseado no tipo de mídia
	userItem := &models.UserItem{
		UserID:       userID,
		ItemID:       itemID,
		Status:       status,
		ProgressType: models.GetDefaultProgressType(item.Type),
		ProgressData: models.JSONB{},
	}

	// Se status for "in_progress", iniciar history
	if status == models.StatusInProgress {
		userItem.StartNewView()
	}

	// Validar
	if err := userItem.Validate(); err != nil {
		return nil, err
	}

	// Criar no banco
	if err := s.userItemRepo.Create(ctx, userItem); err != nil {
		return nil, fmt.Errorf("failed to add item to list: %w", err)
	}

	// Carregar o item completo para retornar
	userItem.Item = *item

	return userItem, nil
}

// GetMyList retorna a lista completa do usuário com paginação
func (s *UserItemService) GetMyList(ctx context.Context, userID uint, params dto.PaginationParams) ([]models.UserItem, int64, error) {
	return s.userItemRepo.GetByUserID(ctx, userID, params)
}

// GetMyListByStatus retorna items da lista filtrados por status com paginação
func (s *UserItemService) GetMyListByStatus(ctx context.Context, userID uint, status models.MediaStatus, params dto.PaginationParams) ([]models.UserItem, int64, error) {
	if !status.IsValid() {
		return nil, 0, models.ErrInvalidStatus
	}
	return s.userItemRepo.GetByStatus(ctx, userID, status, params)
}

// GetMyFavorites retorna todos os items favoritos do usuário com paginação
func (s *UserItemService) GetMyFavorites(ctx context.Context, userID uint, params dto.PaginationParams) ([]models.UserItem, int64, error) {
	return s.userItemRepo.GetFavorites(ctx, userID, params)
}

// GetMyListItem retorna um item específico da lista do usuário
func (s *UserItemService) GetMyListItem(ctx context.Context, id uint, userID uint) (*models.UserItem, error) {
	userItem, err := s.userItemRepo.GetByIDAndUser(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return userItem, nil
}

// UpdateListItem atualiza um item da lista do usuário
func (s *UserItemService) UpdateListItem(ctx context.Context, id uint, userID uint, updates *models.UserItem) (*models.UserItem, error) {
	// Buscar o item existente
	existingItem, err := s.userItemRepo.GetByIDAndUser(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Atualizar campos
	if updates.Status != "" {
		if !updates.Status.IsValid() {
			return nil, models.ErrInvalidStatus
		}

		// Se mudou para "completed", completar visualização atual
		if updates.Status == models.StatusCompleted && existingItem.Status != models.StatusCompleted {
			existingItem.CompleteCurrentView()
		}

		existingItem.Status = updates.Status
	}

	if updates.Rating >= 0 && updates.Rating <= 10 {
		existingItem.Rating = updates.Rating
	} else if updates.Rating < 0 || updates.Rating > 10 {
		return nil, models.ErrInvalidRating
	}

	existingItem.Favorite = updates.Favorite
	existingItem.Notes = updates.Notes

	// Atualizar ProgressType se fornecido
	if updates.ProgressType != "" {
		if !updates.ProgressType.IsValid() {
			return nil, models.ErrInvalidProgressType
		}
		existingItem.ProgressType = updates.ProgressType
	}

	// Atualizar ProgressData se fornecido
	if updates.ProgressData != nil {
		existingItem.ProgressData = updates.ProgressData
	}

	// Atualizar CompletionCount se fornecido
	if updates.CompletionCount > 0 {
		existingItem.CompletionCount = updates.CompletionCount
	}

	// Validar
	if err := existingItem.Validate(); err != nil {
		return nil, err
	}

	// Salvar
	if err := s.userItemRepo.Update(ctx, existingItem); err != nil {
		return nil, fmt.Errorf("failed to update list item: %w", err)
	}

	return existingItem, nil
}

// RemoveFromList remove um item da lista do usuário
func (s *UserItemService) RemoveFromList(ctx context.Context, id uint, userID uint) error {
	// Verificar se o item pertence ao usuário
	_, err := s.userItemRepo.GetByIDAndUser(ctx, id, userID)
	if err != nil {
		return err
	}

	// Remover
	if err := s.userItemRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to remove item from list: %w", err)
	}

	return nil
}

// GetStatistics retorna estatísticas da lista do usuário
func (s *UserItemService) GetStatistics(ctx context.Context, userID uint) (map[string]int64, error) {
	return s.userItemRepo.GetStatistics(ctx, userID)
}
