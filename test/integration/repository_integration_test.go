package integration

import (
	"context"
	"testing"

	"github.com/rafaelc-rb/geekery-api/internal/dto"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/repositories"
	"github.com/rafaelc-rb/geekery-api/internal/testutil"
)

func TestItemRepository_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := repositories.NewItemRepository(db)
	ctx := context.Background()

	t.Run("Create and Get Item", func(t *testing.T) {
		item := &models.Item{
			Title: "Test Anime",
			Type:  models.MediaTypeAnime,
		}

		err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("Failed to create item: %v", err)
		}

		if item.ID == 0 {
			t.Error("Expected item ID to be set")
		}

		retrieved, err := repo.GetByID(ctx, item.ID)
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}

		if retrieved.Title != "Test Anime" {
			t.Errorf("Expected title 'Test Anime', got '%s'", retrieved.Title)
		}
	})

	t.Run("GetAll Items", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		items := []*models.Item{
			{Title: "Item 1", Type: models.MediaTypeAnime},
			{Title: "Item 2", Type: models.MediaTypeMovie},
		}

		for _, item := range items {
			err := repo.Create(ctx, item)
			if err != nil {
				t.Fatalf("Failed to create item: %v", err)
			}
		}

		params := dto.PaginationParams{Page: 1, Limit: 20}
		all, total, err := repo.GetAll(ctx, params)
		if err != nil {
			t.Fatalf("Failed to get all items: %v", err)
		}

		if len(all) != 2 {
			t.Errorf("Expected 2 items, got %d", len(all))
		}
		if total != 2 {
			t.Errorf("Expected total 2, got %d", total)
		}
	})

	t.Run("Update Item", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		item := &models.Item{
			Title: "Original Title",
			Type:  models.MediaTypeAnime,
		}

		err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("Failed to create item: %v", err)
		}

		item.Title = "Updated Title"
		err = repo.Update(ctx, item)
		if err != nil {
			t.Fatalf("Failed to update item: %v", err)
		}

		updated, err := repo.GetByID(ctx, item.ID)
		if err != nil {
			t.Fatalf("Failed to get updated item: %v", err)
		}

		if updated.Title != "Updated Title" {
			t.Errorf("Expected title 'Updated Title', got '%s'", updated.Title)
		}
	})

	t.Run("Delete Item", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		item := &models.Item{
			Title: "To Delete",
			Type:  models.MediaTypeAnime,
		}

		err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("Failed to create item: %v", err)
		}

		err = repo.Delete(ctx, item.ID)
		if err != nil {
			t.Fatalf("Failed to delete item: %v", err)
		}

		_, err = repo.GetByID(ctx, item.ID)
		if err == nil {
			t.Error("Expected error when getting deleted item")
		}
	})

	t.Run("SearchByTitle", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		items := []*models.Item{
			{Title: "Attack on Titan", Type: models.MediaTypeAnime},
			{Title: "Fullmetal Alchemist", Type: models.MediaTypeAnime},
			{Title: "Death Note", Type: models.MediaTypeAnime},
		}

		for _, item := range items {
			err := repo.Create(ctx, item)
			if err != nil {
				t.Fatalf("Failed to create item: %v", err)
			}
		}

		params := dto.PaginationParams{Page: 1, Limit: 20}
		results, total, err := repo.SearchByTitle(ctx, "attack", params)
		if err != nil {
			t.Fatalf("Failed to search: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}
		if total != 1 {
			t.Errorf("Expected total 1, got %d", total)
		}

		if len(results) > 0 && results[0].Title != "Attack on Titan" {
			t.Errorf("Expected 'Attack on Titan', got '%s'", results[0].Title)
		}
	})
}

func TestTagRepository_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := repositories.NewTagRepository(db)
	ctx := context.Background()

	t.Run("Create and Get Tag", func(t *testing.T) {
		tag := &models.Tag{Name: "action"}

		err := repo.Create(ctx, tag)
		if err != nil {
			t.Fatalf("Failed to create tag: %v", err)
		}

		if tag.ID == 0 {
			t.Error("Expected tag ID to be set")
		}

		retrieved, err := repo.GetByID(ctx, tag.ID)
		if err != nil {
			t.Fatalf("Failed to get tag: %v", err)
		}

		if retrieved.Name != "action" {
			t.Errorf("Expected name 'action', got '%s'", retrieved.Name)
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		tag := &models.Tag{Name: "comedy"}
		err := repo.Create(ctx, tag)
		if err != nil {
			t.Fatalf("Failed to create tag: %v", err)
		}

		found, err := repo.GetByName(ctx, "comedy")
		if err != nil {
			t.Fatalf("Failed to get by name: %v", err)
		}

		if found.Name != "comedy" {
			t.Errorf("Expected name 'comedy', got '%s'", found.Name)
		}
	})

	t.Run("FindOrCreate Existing", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		tag := &models.Tag{Name: "drama"}
		err := repo.Create(ctx, tag)
		if err != nil {
			t.Fatalf("Failed to create tag: %v", err)
		}

		newTag := &models.Tag{Name: "drama"}
		err = repo.FindOrCreate(ctx, newTag)
		if err != nil {
			t.Fatalf("Failed to find or create: %v", err)
		}

		if newTag.ID != tag.ID {
			t.Errorf("Expected same ID %d, got %d", tag.ID, newTag.ID)
		}
	})

	t.Run("FindOrCreate New", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		tag := &models.Tag{Name: "horror"}
		err := repo.FindOrCreate(ctx, tag)
		if err != nil {
			t.Fatalf("Failed to find or create: %v", err)
		}

		if tag.ID == 0 {
			t.Error("Expected tag ID to be set")
		}
	})
}

func TestUserItemRepository_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	itemRepo := repositories.NewItemRepository(db)
	userItemRepo := repositories.NewUserItemRepository(db)
	ctx := context.Background()

	// Create test users first
	user1 := &models.User{Name: "testuser1", Email: "test1@example.com"}
	user2 := &models.User{Name: "testuser2", Email: "test2@example.com"}
	db.Create(user1)
	db.Create(user2)

	// Create a test item first
	item := &models.Item{
		Title: "Test Item",
		Type:  models.MediaTypeAnime,
	}
	err := itemRepo.Create(ctx, item)
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	t.Run("Create and Get UserItem", func(t *testing.T) {
		userItem := &models.UserItem{
			UserID:       user1.ID,
			ItemID:       item.ID,
			Status:       models.StatusCompleted,
			Rating:       9.0,
			ProgressType: models.ProgressTypeEpisodic,
		}

		err := userItemRepo.Create(ctx, userItem)
		if err != nil {
			t.Fatalf("Failed to create user item: %v", err)
		}

		if userItem.ID == 0 {
			t.Error("Expected user item ID to be set")
		}

		retrieved, err := userItemRepo.GetByID(ctx, userItem.ID)
		if err != nil {
			t.Fatalf("Failed to get user item: %v", err)
		}

		if retrieved.Status != models.StatusCompleted {
			t.Errorf("Expected status completed, got %s", retrieved.Status)
		}
	})

	t.Run("GetByUserID", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// Recreate items
		item2a := &models.Item{Title: "Item 2a", Type: models.MediaTypeAnime}
		item2b := &models.Item{Title: "Item 2b", Type: models.MediaTypeMovie}
		itemRepo.Create(ctx, item2a)
		itemRepo.Create(ctx, item2b)

		// Recreate user
		user3 := &models.User{Name: "testuser3", Email: "test3@example.com"}
		db.Create(user3)

		userItems := []*models.UserItem{
			{UserID: user3.ID, ItemID: item2a.ID, Status: models.StatusCompleted, ProgressType: models.ProgressTypeEpisodic},
			{UserID: user3.ID, ItemID: item2b.ID, Status: models.StatusInProgress, ProgressType: models.ProgressTypeEpisodic},
		}

		for _, ui := range userItems {
			err := userItemRepo.Create(ctx, ui)
			if err != nil {
				t.Fatalf("Failed to create user item: %v", err)
			}
		}

		params := dto.PaginationParams{Page: 1, Limit: 20}
		all, total, err := userItemRepo.GetByUserID(ctx, user3.ID, params)
		if err != nil {
			t.Fatalf("Failed to get by user ID: %v", err)
		}

		if len(all) != 2 {
			t.Errorf("Expected 2 user items, got %d", len(all))
		}
		if total != 2 {
			t.Errorf("Expected total 2, got %d", total)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		item3 := &models.Item{Title: "Item 3", Type: models.MediaTypeAnime}
		itemRepo.Create(ctx, item3)

		user4 := &models.User{Name: "testuser4", Email: "test4@example.com"}
		db.Create(user4)

		exists, err := userItemRepo.Exists(ctx, user4.ID, item3.ID)
		if err != nil {
			t.Fatalf("Failed to check exists: %v", err)
		}

		if exists {
			t.Error("Expected not to exist")
		}

		userItem := &models.UserItem{
			UserID:       user4.ID,
			ItemID:       item3.ID,
			Status:       models.StatusCompleted,
			ProgressType: models.ProgressTypeEpisodic,
		}
		userItemRepo.Create(ctx, userItem)

		exists, err = userItemRepo.Exists(ctx, user4.ID, item3.ID)
		if err != nil {
			t.Fatalf("Failed to check exists: %v", err)
		}

		if !exists {
			t.Error("Expected to exist")
		}
	})

	t.Run("GetStatistics", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// Create multiple items
		item4a := &models.Item{Title: "Item 4a", Type: models.MediaTypeAnime}
		item4b := &models.Item{Title: "Item 4b", Type: models.MediaTypeMovie}
		item4c := &models.Item{Title: "Item 4c", Type: models.MediaTypeSeries}
		itemRepo.Create(ctx, item4a)
		itemRepo.Create(ctx, item4b)
		itemRepo.Create(ctx, item4c)

		user5 := &models.User{Name: "testuser5", Email: "test5@example.com"}
		db.Create(user5)

		userItems := []*models.UserItem{
			{UserID: user5.ID, ItemID: item4a.ID, Status: models.StatusCompleted, ProgressType: models.ProgressTypeEpisodic},
			{UserID: user5.ID, ItemID: item4b.ID, Status: models.StatusInProgress, ProgressType: models.ProgressTypeEpisodic},
			{UserID: user5.ID, ItemID: item4c.ID, Status: models.StatusPlanned, ProgressType: models.ProgressTypeEpisodic},
		}

		for _, ui := range userItems {
			err := userItemRepo.Create(ctx, ui)
			if err != nil {
				t.Fatalf("Failed to create user item: %v", err)
			}
		}

		stats, err := userItemRepo.GetStatistics(ctx, user5.ID)
		if err != nil {
			t.Fatalf("Failed to get statistics: %v", err)
		}

		if stats["total"] != 3 {
			t.Errorf("Expected total 3, got %d", stats["total"])
		}
	})
}
