package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/services"
	"github.com/rafaelc-rb/geekery-api/internal/testutil"
)

func setupItemHandler() (*ItemHandler, *testutil.MockItemRepository) {
	mockRepo := &testutil.MockItemRepository{}
	service := services.NewItemService(mockRepo)
	handler := NewItemHandler(service)
	gin.SetMode(gin.TestMode)
	return handler, mockRepo
}

func TestItemHandler_GetAllItems(t *testing.T) {
	handler, mockRepo := setupItemHandler()

	expectedItems := []models.Item{
		{Title: "Item 1", Type: models.MediaTypeAnime},
		{Title: "Item 2", Type: models.MediaTypeMovie},
	}

	mockRepo.GetAllFunc = func(ctx context.Context) ([]models.Item, error) {
		return expectedItems, nil
	}

	router := gin.New()
	router.GET("/items", handler.GetAllItems)

	req, _ := http.NewRequest("GET", "/items", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var items []models.Item
	err := json.Unmarshal(w.Body.Bytes(), &items)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
}

func TestItemHandler_GetItemByID(t *testing.T) {
	handler, mockRepo := setupItemHandler()

	expectedItem := &models.Item{Title: "Test Item", Type: models.MediaTypeAnime}
	expectedItem.ID = 1

	mockRepo.GetByIDFunc = func(ctx context.Context, id uint) (*models.Item, error) {
		if id == 1 {
			return expectedItem, nil
		}
		return nil, nil
	}

	router := gin.New()
	router.GET("/items/:id", handler.GetItemByID)

	req, _ := http.NewRequest("GET", "/items/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var item models.Item
	err := json.Unmarshal(w.Body.Bytes(), &item)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if item.Title != "Test Item" {
		t.Errorf("Expected title 'Test Item', got '%s'", item.Title)
	}
}

func TestItemHandler_GetItemByID_Invalid(t *testing.T) {
	handler, _ := setupItemHandler()

	router := gin.New()
	router.GET("/items/:id", handler.GetItemByID)

	req, _ := http.NewRequest("GET", "/items/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestItemHandler_CreateItem(t *testing.T) {
	handler, mockRepo := setupItemHandler()

	mockRepo.CreateFunc = func(ctx context.Context, item *models.Item) error {
		item.ID = 1
		return nil
	}
	mockRepo.AssociateTagsFunc = func(ctx context.Context, itemID uint, tagIDs []uint) error {
		return nil
	}

	router := gin.New()
	router.POST("/items", handler.CreateItem)

	payload := map[string]interface{}{
		"title": "New Item",
		"type":  "anime",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestItemHandler_UpdateItem(t *testing.T) {
	handler, mockRepo := setupItemHandler()

	existingItem := &models.Item{Title: "Old Title", Type: models.MediaTypeAnime}
	existingItem.ID = 1

	mockRepo.GetByIDFunc = func(ctx context.Context, id uint) (*models.Item, error) {
		return existingItem, nil
	}
	mockRepo.UpdateFunc = func(ctx context.Context, item *models.Item) error {
		return nil
	}
	mockRepo.AssociateTagsFunc = func(ctx context.Context, itemID uint, tagIDs []uint) error {
		return nil
	}

	router := gin.New()
	router.PUT("/items/:id", handler.UpdateItem)

	payload := map[string]interface{}{
		"title": "Updated Title",
		"type":  "anime",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/items/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestItemHandler_DeleteItem(t *testing.T) {
	handler, mockRepo := setupItemHandler()

	mockRepo.DeleteFunc = func(ctx context.Context, id uint) error {
		return nil
	}

	router := gin.New()
	router.DELETE("/items/:id", handler.DeleteItem)

	req, _ := http.NewRequest("DELETE", "/items/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestItemHandler_SearchItems(t *testing.T) {
	handler, mockRepo := setupItemHandler()

	expectedItems := []models.Item{
		{Title: "Attack on Titan", Type: models.MediaTypeAnime},
	}

	mockRepo.SearchByTitleFunc = func(ctx context.Context, query string) ([]models.Item, error) {
		return expectedItems, nil
	}

	router := gin.New()
	router.GET("/items/search", handler.SearchItems)

	req, _ := http.NewRequest("GET", "/items/search?q=attack", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var items []models.Item
	err := json.Unmarshal(w.Body.Bytes(), &items)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}
}

func TestItemHandler_SearchItems_MissingQuery(t *testing.T) {
	handler, mockRepo := setupItemHandler()

	// O handler não valida query vazio, então retorna array vazio com status 200
	mockRepo.SearchByTitleFunc = func(ctx context.Context, query string) ([]models.Item, error) {
		return []models.Item{}, nil
	}

	router := gin.New()
	router.GET("/items/search", handler.SearchItems)

	req, _ := http.NewRequest("GET", "/items/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
