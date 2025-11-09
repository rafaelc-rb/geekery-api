package services

import (
	"context"
	"errors"
	"testing"

	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/testutil"
)

func TestCreateItem_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &testutil.MockItemRepository{
		CreateFunc: func(ctx context.Context, item *models.Item) error {
			item.ID = 1
			return nil
		},
		AssociateTagsFunc: func(ctx context.Context, itemID uint, tagIDs []uint) error {
			return nil
		},
	}

	service := NewItemService(mockRepo)

	item := &models.Item{
		Title: "Test Item",
		Type:  models.MediaTypeAnime,
	}

	err := service.CreateItem(ctx, item, []uint{1, 2})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if item.ID != 1 {
		t.Errorf("Expected item ID to be set to 1, got %d", item.ID)
	}
}

func TestCreateItem_ValidationError(t *testing.T) {
	ctx := context.Background()
	mockRepo := &testutil.MockItemRepository{}
	service := NewItemService(mockRepo)

	item := &models.Item{
		Title: "", // Invalid: empty title
		Type:  models.MediaTypeAnime,
	}

	err := service.CreateItem(ctx, item, nil)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
}

func TestCreateItemWithSpecificData_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &testutil.MockItemRepository{
		CreateFunc: func(ctx context.Context, item *models.Item) error {
			item.ID = 1
			return nil
		},
		CreateSpecificDataFunc: func(ctx context.Context, itemID uint, mediaType models.MediaType, data interface{}) error {
			return nil
		},
		AssociateTagsFunc: func(ctx context.Context, itemID uint, tagIDs []uint) error {
			return nil
		},
	}

	service := NewItemService(mockRepo)

	item := &models.Item{
		Title: "Attack on Titan",
		Type:  models.MediaTypeAnime,
	}

	animeData := &models.AnimeData{
		Episodes: 75,
		Studio:   "MAPPA",
	}

	err := service.CreateItemWithSpecificData(ctx, item, animeData, []uint{1})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetAllItems_Success(t *testing.T) {
	ctx := context.Background()
	expectedItems := []models.Item{
		{Title: "Item 1", Type: models.MediaTypeAnime},
		{Title: "Item 2", Type: models.MediaTypeMovie},
	}

	mockRepo := &testutil.MockItemRepository{
		GetAllFunc: func(ctx context.Context) ([]models.Item, error) {
			return expectedItems, nil
		},
	}

	service := NewItemService(mockRepo)
	items, err := service.GetAllItems(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
}

func TestGetItemByID_Success(t *testing.T) {
	ctx := context.Background()
	expectedItem := &models.Item{
		Title: "Test Item",
		Type:  models.MediaTypeAnime,
	}
	expectedItem.ID = 1

	mockRepo := &testutil.MockItemRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*models.Item, error) {
			return expectedItem, nil
		},
	}

	service := NewItemService(mockRepo)
	item, err := service.GetItemByID(ctx, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if item.Title != "Test Item" {
		t.Errorf("Expected title 'Test Item', got '%s'", item.Title)
	}
}

func TestGetItemByID_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := &testutil.MockItemRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*models.Item, error) {
			return nil, errors.New("record not found")
		},
	}

	service := NewItemService(mockRepo)
	_, err := service.GetItemByID(ctx, 999)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestUpdateItem_Success(t *testing.T) {
	ctx := context.Background()
	existingItem := &models.Item{
		Title: "Old Title",
		Type:  models.MediaTypeAnime,
	}
	existingItem.ID = 1

	mockRepo := &testutil.MockItemRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*models.Item, error) {
			return existingItem, nil
		},
		UpdateFunc: func(ctx context.Context, item *models.Item) error {
			return nil
		},
		AssociateTagsFunc: func(ctx context.Context, itemID uint, tagIDs []uint) error {
			return nil
		},
	}

	service := NewItemService(mockRepo)

	updatedItem := &models.Item{
		Title: "New Title",
		Type:  models.MediaTypeAnime,
	}

	err := service.UpdateItem(ctx, 1, updatedItem, []uint{1, 2})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteItem_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &testutil.MockItemRepository{
		DeleteFunc: func(ctx context.Context, id uint) error {
			return nil
		},
	}

	service := NewItemService(mockRepo)
	err := service.DeleteItem(ctx, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
