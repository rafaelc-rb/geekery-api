package services

import (
	"context"
	"testing"

	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/testutil"
	"gorm.io/gorm"
)

func TestCreateTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &testutil.MockTagRepository{
		GetByNameFunc: func(ctx context.Context, name string) (*models.Tag, error) {
			return nil, gorm.ErrRecordNotFound
		},
		CreateFunc: func(ctx context.Context, tag *models.Tag) error {
			tag.ID = 1
			return nil
		},
	}

	service := NewTagService(mockRepo)
	tag := &models.Tag{Name: "Action"}

	err := service.CreateTag(ctx, tag)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if tag.ID != 1 {
		t.Errorf("Expected tag ID to be set to 1, got %d", tag.ID)
	}
}

func TestCreateTag_EmptyName(t *testing.T) {
	ctx := context.Background()
	mockRepo := &testutil.MockTagRepository{}
	service := NewTagService(mockRepo)

	tag := &models.Tag{Name: ""}

	err := service.CreateTag(ctx, tag)

	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}
}

func TestCreateTag_DuplicateName(t *testing.T) {
	ctx := context.Background()
	existingTag := &models.Tag{Name: "action"}
	existingTag.ID = 1

	mockRepo := &testutil.MockTagRepository{
		GetByNameFunc: func(ctx context.Context, name string) (*models.Tag, error) {
			return existingTag, nil
		},
	}

	service := NewTagService(mockRepo)
	tag := &models.Tag{Name: "Action"} // Different case, but should be caught

	err := service.CreateTag(ctx, tag)

	if err != models.ErrDuplicateTag {
		t.Errorf("Expected duplicate tag error, got %v", err)
	}
}

func TestGetAllTags_Success(t *testing.T) {
	ctx := context.Background()
	expectedTags := []models.Tag{
		{Name: "Action"},
		{Name: "Comedy"},
	}

	mockRepo := &testutil.MockTagRepository{
		GetAllFunc: func(ctx context.Context) ([]models.Tag, error) {
			return expectedTags, nil
		},
	}

	service := NewTagService(mockRepo)
	tags, err := service.GetAllTags(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
}

func TestGetTagByID_Success(t *testing.T) {
	ctx := context.Background()
	expectedTag := &models.Tag{Name: "Action"}
	expectedTag.ID = 1

	mockRepo := &testutil.MockTagRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*models.Tag, error) {
			return expectedTag, nil
		},
	}

	service := NewTagService(mockRepo)
	tag, err := service.GetTagByID(ctx, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if tag.Name != "Action" {
		t.Errorf("Expected name 'Action', got '%s'", tag.Name)
	}
}

func TestGetTagByID_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := &testutil.MockTagRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*models.Tag, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := NewTagService(mockRepo)
	_, err := service.GetTagByID(ctx, 999)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestUpdateTag_Success(t *testing.T) {
	ctx := context.Background()
	existingTag := &models.Tag{Name: "Old Name"}
	existingTag.ID = 1

	mockRepo := &testutil.MockTagRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*models.Tag, error) {
			return existingTag, nil
		},
		GetByNameFunc: func(ctx context.Context, name string) (*models.Tag, error) {
			return nil, gorm.ErrRecordNotFound
		},
		UpdateFunc: func(ctx context.Context, tag *models.Tag) error {
			return nil
		},
	}

	service := NewTagService(mockRepo)

	updatedTag := &models.Tag{Name: "New Name"}

	err := service.UpdateTag(ctx, 1, updatedTag)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &testutil.MockTagRepository{
		DeleteFunc: func(ctx context.Context, id uint) error {
			return nil
		},
	}

	service := NewTagService(mockRepo)
	err := service.DeleteTag(ctx, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestFindOrCreateTags_AllExist(t *testing.T) {
	ctx := context.Background()
	existingTags := []models.Tag{
		{Name: "action"},
		{Name: "comedy"},
	}
	existingTags[0].ID = 1
	existingTags[1].ID = 2

	mockRepo := &testutil.MockTagRepository{
		FindOrCreateFunc: func(ctx context.Context, tag *models.Tag) error {
			switch tag.Name {
			case "action":
				tag.ID = 1
			case "comedy":
				tag.ID = 2
			}
			return nil
		},
	}

	service := NewTagService(mockRepo)
	tagNames := []string{"Action", "Comedy"}

	tags, err := service.FindOrCreateTags(ctx, tagNames)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
}

func TestFindOrCreateTags_EmptyNames(t *testing.T) {
	ctx := context.Background()
	mockRepo := &testutil.MockTagRepository{}
	service := NewTagService(mockRepo)

	tagNames := []string{"", "  "}

	tags, err := service.FindOrCreateTags(ctx, tagNames)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(tags))
	}
}

func TestFindOrCreateTags_CreateNew(t *testing.T) {
	ctx := context.Background()
	callCount := 0

	mockRepo := &testutil.MockTagRepository{
		FindOrCreateFunc: func(ctx context.Context, tag *models.Tag) error {
			callCount++
			tag.ID = uint(callCount)
			return nil
		},
	}

	service := NewTagService(mockRepo)
	tagNames := []string{"NewTag1", "NewTag2"}

	tags, err := service.FindOrCreateTags(ctx, tagNames)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
	if callCount != 2 {
		t.Errorf("Expected FindOrCreate to be called 2 times, got %d", callCount)
	}
}
