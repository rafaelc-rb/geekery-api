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
	"gorm.io/gorm"
)

func setupTagHandler() (*TagHandler, *testutil.MockTagRepository) {
	mockRepo := &testutil.MockTagRepository{}
	service := services.NewTagService(mockRepo)
	handler := NewTagHandler(service)
	gin.SetMode(gin.TestMode)
	return handler, mockRepo
}

func TestTagHandler_GetAllTags(t *testing.T) {
	handler, mockRepo := setupTagHandler()

	expectedTags := []models.Tag{
		{Name: "action"},
		{Name: "comedy"},
	}

	mockRepo.GetAllFunc = func(ctx context.Context) ([]models.Tag, error) {
		return expectedTags, nil
	}

	router := gin.New()
	router.GET("/tags", handler.GetAllTags)

	req, _ := http.NewRequest("GET", "/tags", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var tags []models.Tag
	err := json.Unmarshal(w.Body.Bytes(), &tags)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
}

func TestTagHandler_GetTagByID(t *testing.T) {
	handler, mockRepo := setupTagHandler()

	expectedTag := &models.Tag{Name: "action"}
	expectedTag.ID = 1

	mockRepo.GetByIDFunc = func(ctx context.Context, id uint) (*models.Tag, error) {
		if id == 1 {
			return expectedTag, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	router := gin.New()
	router.GET("/tags/:id", handler.GetTagByID)

	req, _ := http.NewRequest("GET", "/tags/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var tag models.Tag
	err := json.Unmarshal(w.Body.Bytes(), &tag)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if tag.Name != "action" {
		t.Errorf("Expected name 'action', got '%s'", tag.Name)
	}
}

func TestTagHandler_GetTagByID_Invalid(t *testing.T) {
	handler, _ := setupTagHandler()

	router := gin.New()
	router.GET("/tags/:id", handler.GetTagByID)

	req, _ := http.NewRequest("GET", "/tags/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestTagHandler_CreateTag(t *testing.T) {
	handler, mockRepo := setupTagHandler()

	mockRepo.GetByNameFunc = func(ctx context.Context, name string) (*models.Tag, error) {
		return nil, gorm.ErrRecordNotFound
	}
	mockRepo.CreateFunc = func(ctx context.Context, tag *models.Tag) error {
		tag.ID = 1
		return nil
	}

	router := gin.New()
	router.POST("/tags", handler.CreateTag)

	payload := map[string]interface{}{
		"name": "NewTag",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestTagHandler_CreateTag_EmptyName(t *testing.T) {
	handler, _ := setupTagHandler()

	router := gin.New()
	router.POST("/tags", handler.CreateTag)

	payload := map[string]interface{}{
		"name": "",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestTagHandler_UpdateTag(t *testing.T) {
	handler, mockRepo := setupTagHandler()

	existingTag := &models.Tag{Name: "oldname"}
	existingTag.ID = 1

	mockRepo.GetByIDFunc = func(ctx context.Context, id uint) (*models.Tag, error) {
		return existingTag, nil
	}
	mockRepo.GetByNameFunc = func(ctx context.Context, name string) (*models.Tag, error) {
		return nil, gorm.ErrRecordNotFound
	}
	mockRepo.UpdateFunc = func(ctx context.Context, tag *models.Tag) error {
		return nil
	}

	router := gin.New()
	router.PUT("/tags/:id", handler.UpdateTag)

	payload := map[string]interface{}{
		"name": "newname",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/tags/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestTagHandler_DeleteTag(t *testing.T) {
	handler, mockRepo := setupTagHandler()

	mockRepo.DeleteFunc = func(ctx context.Context, id uint) error {
		return nil
	}

	router := gin.New()
	router.DELETE("/tags/:id", handler.DeleteTag)

	req, _ := http.NewRequest("DELETE", "/tags/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
