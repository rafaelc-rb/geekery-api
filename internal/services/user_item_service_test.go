package services

import (
	"context"
	"testing"

	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/testutil"
	"gorm.io/gorm"
)

func TestAddToList_Success(t *testing.T) {
	ctx := context.Background()
	item := &models.Item{Title: "Test Item", Type: models.MediaTypeAnime}
	item.ID = 1

	mockItemRepo := &testutil.MockItemRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*models.Item, error) {
			return item, nil
		},
	}

	mockUserItemRepo := &testutil.MockUserItemRepository{
		ExistsFunc: func(ctx context.Context, userID, itemID uint) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, userItem *models.UserItem) error {
			userItem.ID = 1
			return nil
		},
	}

	service := NewUserItemService(mockUserItemRepo, mockItemRepo)
	userItem, err := service.AddToList(ctx, 1, 1, models.StatusPlanned)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if userItem.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", userItem.UserID)
	}
	if userItem.Status != models.StatusPlanned {
		t.Errorf("Expected status planned, got %s", userItem.Status)
	}
}

func TestAddToList_ItemNotFound(t *testing.T) {
	ctx := context.Background()
	mockItemRepo := &testutil.MockItemRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*models.Item, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	mockUserItemRepo := &testutil.MockUserItemRepository{}
	service := NewUserItemService(mockUserItemRepo, mockItemRepo)

	_, err := service.AddToList(ctx, 1, 999, models.StatusPlanned)

	if err == nil {
		t.Error("Expected error for item not found, got nil")
	}
}

func TestAddToList_DuplicateEntry(t *testing.T) {
	ctx := context.Background()
	item := &models.Item{Title: "Test Item", Type: models.MediaTypeAnime}
	item.ID = 1

	mockItemRepo := &testutil.MockItemRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*models.Item, error) {
			return item, nil
		},
	}

	mockUserItemRepo := &testutil.MockUserItemRepository{
		ExistsFunc: func(ctx context.Context, userID, itemID uint) (bool, error) {
			return true, nil // Already exists
		},
	}

	service := NewUserItemService(mockUserItemRepo, mockItemRepo)
	_, err := service.AddToList(ctx, 1, 1, models.StatusPlanned)

	if err != models.ErrDuplicateEntry {
		t.Errorf("Expected duplicate entry error, got %v", err)
	}
}

func TestGetMyList_Success(t *testing.T) {
	ctx := context.Background()
	expectedItems := []models.UserItem{
		{UserID: 1, ItemID: 1, Status: models.StatusCompleted},
		{UserID: 1, ItemID: 2, Status: models.StatusInProgress},
	}

	mockUserItemRepo := &testutil.MockUserItemRepository{
		GetByUserIDFunc: func(ctx context.Context, userID uint) ([]models.UserItem, error) {
			return expectedItems, nil
		},
	}

	service := NewUserItemService(mockUserItemRepo, nil)
	items, err := service.GetMyList(ctx, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
}

func TestGetMyListByStatus_Success(t *testing.T) {
	ctx := context.Background()
	expectedItems := []models.UserItem{
		{UserID: 1, ItemID: 1, Status: models.StatusCompleted},
	}

	mockUserItemRepo := &testutil.MockUserItemRepository{
		GetByStatusFunc: func(ctx context.Context, userID uint, status models.MediaStatus) ([]models.UserItem, error) {
			return expectedItems, nil
		},
	}

	service := NewUserItemService(mockUserItemRepo, nil)
	items, err := service.GetMyListByStatus(ctx, 1, models.StatusCompleted)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}
}

func TestUpdateListItem_Success(t *testing.T) {
	ctx := context.Background()
	existingItem := &models.UserItem{
		UserID:       1,
		ItemID:       1,
		Status:       models.StatusInProgress,
		Rating:       7.0,
		ProgressType: models.ProgressTypeEpisodic,
	}
	existingItem.ID = 1

	mockUserItemRepo := &testutil.MockUserItemRepository{
		GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*models.UserItem, error) {
			return existingItem, nil
		},
		UpdateFunc: func(ctx context.Context, userItem *models.UserItem) error {
			return nil
		},
	}

	service := NewUserItemService(mockUserItemRepo, nil)

	updates := &models.UserItem{
		Status: models.StatusCompleted,
		Rating: 9.0,
	}

	updated, err := service.UpdateListItem(ctx, 1, 1, updates)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updated.Status != models.StatusCompleted {
		t.Errorf("Expected status completed, got %s", updated.Status)
	}
	if updated.Rating != 9.0 {
		t.Errorf("Expected rating 9.0, got %.1f", updated.Rating)
	}
}

func TestUpdateListItem_CompleteView(t *testing.T) {
	ctx := context.Background()
	existingItem := &models.UserItem{
		UserID:          1,
		ItemID:          1,
		Status:          models.StatusInProgress,
		ProgressType:    models.ProgressTypeEpisodic,
		CompletionCount: 0,
	}
	existingItem.ID = 1
	existingItem.StartNewView() // Start a view

	mockUserItemRepo := &testutil.MockUserItemRepository{
		GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*models.UserItem, error) {
			return existingItem, nil
		},
		UpdateFunc: func(ctx context.Context, userItem *models.UserItem) error {
			return nil
		},
	}

	service := NewUserItemService(mockUserItemRepo, nil)

	updates := &models.UserItem{
		Status: models.StatusCompleted,
	}

	updated, err := service.UpdateListItem(ctx, 1, 1, updates)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updated.CompletionCount != 1 {
		t.Errorf("Expected CompletionCount 1, got %d", updated.CompletionCount)
	}
}

func TestRemoveFromList_Success(t *testing.T) {
	ctx := context.Background()
	existingItem := &models.UserItem{UserID: 1, ItemID: 1}
	existingItem.ID = 1

	mockUserItemRepo := &testutil.MockUserItemRepository{
		GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*models.UserItem, error) {
			return existingItem, nil
		},
		DeleteFunc: func(ctx context.Context, id uint) error {
			return nil
		},
	}

	service := NewUserItemService(mockUserItemRepo, nil)
	err := service.RemoveFromList(ctx, 1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetStatistics_Success(t *testing.T) {
	ctx := context.Background()
	expectedStats := map[string]int64{
		"total":       10,
		"completed":   5,
		"in_progress": 3,
		"planned":     2,
	}

	mockUserItemRepo := &testutil.MockUserItemRepository{
		GetStatisticsFunc: func(ctx context.Context, userID uint) (map[string]int64, error) {
			return expectedStats, nil
		},
	}

	service := NewUserItemService(mockUserItemRepo, nil)
	stats, err := service.GetStatistics(ctx, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stats["total"] != 10 {
		t.Errorf("Expected total 10, got %d", stats["total"])
	}
}

func TestGetMyFavorites_Success(t *testing.T) {
	ctx := context.Background()
	expectedItems := []models.UserItem{
		{UserID: 1, ItemID: 1, Favorite: true},
		{UserID: 1, ItemID: 2, Favorite: true},
	}

	mockUserItemRepo := &testutil.MockUserItemRepository{
		GetFavoritesFunc: func(ctx context.Context, userID uint) ([]models.UserItem, error) {
			return expectedItems, nil
		},
	}

	service := NewUserItemService(mockUserItemRepo, nil)
	items, err := service.GetMyFavorites(ctx, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 favorites, got %d", len(items))
	}
}

func TestUpdateListItem_InvalidRating(t *testing.T) {
	ctx := context.Background()
	existingItem := &models.UserItem{
		UserID: 1,
		ItemID: 1,
		Status: models.StatusInProgress,
	}
	existingItem.ID = 1

	mockUserItemRepo := &testutil.MockUserItemRepository{
		GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*models.UserItem, error) {
			return existingItem, nil
		},
	}

	service := NewUserItemService(mockUserItemRepo, nil)

	updates := &models.UserItem{
		Rating: 11.0, // Invalid rating
	}

	_, err := service.UpdateListItem(ctx, 1, 1, updates)

	if err != models.ErrInvalidRating {
		t.Errorf("Expected invalid rating error, got %v", err)
	}
}
