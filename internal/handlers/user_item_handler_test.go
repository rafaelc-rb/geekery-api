package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/dto"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/services"
	"github.com/rafaelc-rb/geekery-api/internal/testutil"
)

func setupUserItemHandler() (*UserItemHandler, *testutil.MockUserItemRepository, *testutil.MockItemRepository) {
	mockUserItemRepo := &testutil.MockUserItemRepository{}
	mockItemRepo := &testutil.MockItemRepository{}
	service := services.NewUserItemService(mockUserItemRepo, mockItemRepo)
	handler := NewUserItemHandler(service)
	gin.SetMode(gin.TestMode)
	return handler, mockUserItemRepo, mockItemRepo
}

// mockAuthMiddleware injects a mock userID into the context for testing
func mockAuthMiddleware(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

func TestUserItemHandler_AddToList(t *testing.T) {
	handler, mockUserItemRepo, mockItemRepo := setupUserItemHandler()

	item := &models.Item{Title: "Test Item", Type: models.MediaTypeAnime}
	item.ID = 1

	mockItemRepo.GetByIDFunc = func(ctx context.Context, id uint) (*models.Item, error) {
		return item, nil
	}
	mockUserItemRepo.ExistsFunc = func(ctx context.Context, userID, itemID uint) (bool, error) {
		return false, nil
	}
	mockUserItemRepo.CreateFunc = func(ctx context.Context, userItem *models.UserItem) error {
		userItem.ID = 1
		return nil
	}

	router := gin.New()
	router.Use(mockAuthMiddleware(1)) // Mock authenticated user with ID 1
	router.POST("/my-list", handler.AddToList)

	payload := map[string]interface{}{
		"item_id": 1,
		"status":  "planned",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/my-list", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestUserItemHandler_GetMyList(t *testing.T) {
	handler, mockUserItemRepo, _ := setupUserItemHandler()

	expectedItems := []models.UserItem{
		{UserID: 1, ItemID: 1, Status: models.StatusCompleted},
		{UserID: 1, ItemID: 2, Status: models.StatusInProgress},
	}

	mockUserItemRepo.GetByUserIDFunc = func(ctx context.Context, userID uint, params dto.PaginationParams) ([]models.UserItem, int64, error) {
		return expectedItems, 2, nil
	}

	router := gin.New()
	router.Use(mockAuthMiddleware(1)) // Mock authenticated user with ID 1
	router.GET("/my-list", handler.GetMyList)

	req, _ := http.NewRequest("GET", "/my-list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response dto.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Convert data back to items
	itemsData, _ := json.Marshal(response.Data)
	var items []models.UserItem
	if err := json.Unmarshal(itemsData, &items); err != nil {
		t.Fatalf("Failed to unmarshal items: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
}

func TestUserItemHandler_UpdateListItem(t *testing.T) {
	handler, mockUserItemRepo, _ := setupUserItemHandler()

	existingItem := &models.UserItem{
		UserID:       1,
		ItemID:       1,
		Status:       models.StatusInProgress,
		Rating:       7.0,
		ProgressType: models.ProgressTypeEpisodic,
	}
	existingItem.ID = 1

	mockUserItemRepo.GetByIDAndUserFunc = func(ctx context.Context, id, userID uint) (*models.UserItem, error) {
		return existingItem, nil
	}
	mockUserItemRepo.UpdateFunc = func(ctx context.Context, userItem *models.UserItem) error {
		return nil
	}

	router := gin.New()
	router.Use(mockAuthMiddleware(1)) // Mock authenticated user with ID 1
	router.PUT("/my-list/:id", handler.UpdateListItem)

	payload := map[string]interface{}{
		"status": "completed",
		"rating": 9.0,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/my-list/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestUserItemHandler_RemoveFromList(t *testing.T) {
	handler, mockUserItemRepo, _ := setupUserItemHandler()

	existingItem := &models.UserItem{UserID: 1, ItemID: 1}
	existingItem.ID = 1

	mockUserItemRepo.GetByIDAndUserFunc = func(ctx context.Context, id, userID uint) (*models.UserItem, error) {
		return existingItem, nil
	}
	mockUserItemRepo.DeleteFunc = func(ctx context.Context, id uint) error {
		return nil
	}

	router := gin.New()
	router.Use(mockAuthMiddleware(1)) // Mock authenticated user with ID 1
	router.DELETE("/my-list/:id", handler.RemoveFromList)

	req, _ := http.NewRequest("DELETE", "/my-list/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestUserItemHandler_GetStatistics(t *testing.T) {
	handler, mockUserItemRepo, _ := setupUserItemHandler()

	expectedStats := map[string]int64{
		"total":       10,
		"completed":   5,
		"in_progress": 3,
		"planned":     2,
	}

	mockUserItemRepo.GetStatisticsFunc = func(ctx context.Context, userID uint) (map[string]int64, error) {
		return expectedStats, nil
	}

	router := gin.New()
	router.Use(mockAuthMiddleware(1)) // Mock authenticated user with ID 1
	router.GET("/my-list/stats", handler.GetStatistics)

	req, _ := http.NewRequest("GET", "/my-list/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var stats map[string]int64
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if stats["total"] != 10 {
		t.Errorf("Expected total 10, got %d", stats["total"])
	}
}
