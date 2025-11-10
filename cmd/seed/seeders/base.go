package seeders

import (
	"fmt"
	"time"

	"github.com/rafaelc-rb/geekery-api/internal/models"
	"gorm.io/gorm"
)

// ParseDate helper para converter string em *time.Time
func ParseDate(dateStr string) *time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	return &t
}

// CreateOrSkip cria um registro ou pula se já existe
// Retorna true se foi criado, false se já existia
func CreateOrSkip(db *gorm.DB, model interface{}, where interface{}) (bool, error) {
	if err := db.Where(where).FirstOrCreate(model).Error; err != nil {
		return false, err
	}
	return db.RowsAffected > 0, nil
}

// ItemWithData representa um item do catálogo com seus dados específicos
type ItemWithData struct {
	Item         models.Item
	SpecificData interface{}
	TagNames     []string // Nomes das tags para facilitar associação
}

// CreateItemWithSpecificData cria um item e seus dados específicos
func CreateItemWithSpecificData(db *gorm.DB, itemData ItemWithData, allTags []models.Tag) error {
	// Check if item already exists by title
	var existing models.Item
	if err := db.Where("title = ?", itemData.Item.Title).First(&existing).Error; err == nil {
		fmt.Printf("⚠ Item already exists: %s (ID: %d)\n", existing.Title, existing.ID)
		return nil
	}

	// Create item
	if err := db.Create(&itemData.Item).Error; err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	// Create specific data
	if err := createSpecificData(db, &itemData.Item, itemData.SpecificData); err != nil {
		return err
	}

	// Associate tags
	if len(itemData.TagNames) > 0 {
		var itemTags []models.Tag
		for _, tagName := range itemData.TagNames {
			for _, tag := range allTags {
				if tag.Name == tagName {
					itemTags = append(itemTags, tag)
					break
				}
			}
		}
		if err := db.Model(&itemData.Item).Association("Tags").Replace(itemTags); err != nil {
			return fmt.Errorf("failed to associate tags: %w", err)
		}
	}

	fmt.Printf("✓ Item created: %s (ID: %d, Type: %s)\n", itemData.Item.Title, itemData.Item.ID, itemData.Item.Type)
	return nil
}

// createSpecificData cria os dados específicos baseado no tipo do item
func createSpecificData(db *gorm.DB, item *models.Item, specificData interface{}) error {
	if specificData == nil {
		return nil
	}

	switch item.Type {
	case models.MediaTypeAnime:
		animeData := specificData.(*models.AnimeData)
		animeData.ItemID = item.ID
		return db.Create(animeData).Error
	case models.MediaTypeMovie:
		movieData := specificData.(*models.MovieData)
		movieData.ItemID = item.ID
		return db.Create(movieData).Error
	case models.MediaTypeGame:
		gameData := specificData.(*models.GameData)
		gameData.ItemID = item.ID
		return db.Create(gameData).Error
	case models.MediaTypeSeries:
		seriesData := specificData.(*models.SeriesData)
		seriesData.ItemID = item.ID
		return db.Create(seriesData).Error
	case models.MediaTypeComic, models.MediaTypeNovel, models.MediaTypeBook:
		bookData := specificData.(*models.BookData)
		bookData.ItemID = item.ID
		return db.Create(bookData).Error
	}

	return nil
}

// CreateUserItemIfNotExists cria um user item se não existir
func CreateUserItemIfNotExists(db *gorm.DB, userItem *models.UserItem) error {
	var existing models.UserItem
	if err := db.Where("user_id = ? AND item_id = ?", userItem.UserID, userItem.ItemID).First(&existing).Error; err == nil {
		fmt.Printf("⚠ User item already exists: User %d - Item %d\n", existing.UserID, existing.ItemID)
		return nil
	}

	if err := db.Create(userItem).Error; err != nil {
		return fmt.Errorf("failed to create user item: %w", err)
	}

	fmt.Printf("✓ User item created: User %d - Item %d (Status: %s)\n", userItem.UserID, userItem.ItemID, userItem.Status)
	return nil
}
