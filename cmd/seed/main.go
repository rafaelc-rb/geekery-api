package main

import (
	"fmt"
	"log"

	"github.com/rafaelc-rb/geekery-api/cmd/seed/data"
	"github.com/rafaelc-rb/geekery-api/cmd/seed/seeders"
	"github.com/rafaelc-rb/geekery-api/internal/config"
	"github.com/rafaelc-rb/geekery-api/internal/database"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("🌱 Starting database seeding...")
	fmt.Println()

	// Carregar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Conectar ao banco de dados
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// ========================================
	// 1. Seed Users
	// ========================================
	fmt.Println("👤 Seeding users...")
	user := seedUser(db)
	fmt.Println()

	// ========================================
	// 2. Seed Tags
	// ========================================
	fmt.Println("🏷️  Seeding tags...")
	tags := seedTags(db)
	fmt.Println()

	// ========================================
	// 3. Seed Catalog Items
	// ========================================
	fmt.Println("📚 Seeding catalog items...")
	itemIDs := seedCatalogItems(db, tags)
	fmt.Println()

	// ========================================
	// 4. Seed User Items (Personal Lists)
	// ========================================
	fmt.Println("📝 Seeding user items...")
	seedUserItems(db, user.ID, itemIDs)
	fmt.Println()

	// ========================================
	// Summary
	// ========================================
	printSummary(len(tags), len(itemIDs))
}

// seedUser cria o usuário de exemplo
func seedUser(db *gorm.DB) models.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Failed to hash password: %v", err)
	}

	user := models.User{
		Email:        "demo@geekery.com",
		Username:     "demo",
		Name:         "Demo User",
		PasswordHash: string(hashedPassword),
	}

	if err := db.Where("email = ?", user.Email).FirstOrCreate(&user).Error; err != nil {
		log.Fatalf("❌ Failed to create user: %v", err)
	}

	fmt.Printf("  ✓ User: %s (ID: %d)\n", user.Email, user.ID)
	return user
}

// seedTags cria as tags
func seedTags(db *gorm.DB) []models.Tag {
	tags := data.GetTags()
	createdCount := 0

	for i := range tags {
		created, err := seeders.CreateOrSkip(db, &tags[i], models.Tag{Name: tags[i].Name})
		if err != nil {
			log.Fatalf("❌ Failed to create tag: %v", err)
		}
		if created {
			createdCount++
		}
	}

	fmt.Printf("  ✓ Created %d new tags (total: %d)\n", createdCount, len(tags))
	return tags
}

// seedCatalogItems cria os items do catálogo e retorna um map de título -> ID
func seedCatalogItems(db *gorm.DB, tags []models.Tag) map[string]uint {
	items := data.GetCatalogItems()
	itemIDs := make(map[string]uint)
	createdCount := 0

	for _, itemData := range items {
		// Check if already exists
		var existing models.Item
		if err := db.Where("title = ?", itemData.Item.Title).First(&existing).Error; err == nil {
			itemIDs[existing.Title] = existing.ID
			continue
		}

		// Create item
		if err := seeders.CreateItemWithSpecificData(db, itemData, tags); err != nil {
			log.Fatalf("❌ Failed to create item: %v", err)
		}

		itemIDs[itemData.Item.Title] = itemData.Item.ID
		createdCount++
	}

	fmt.Printf("  ✓ Created %d new items (total: %d)\n", createdCount, len(items))
	return itemIDs
}

// seedUserItems cria os items na lista pessoal do usuário
func seedUserItems(db *gorm.DB, userID uint, itemIDs map[string]uint) {
	userItems := data.GetUserItems(userID, itemIDs)
	createdCount := 0

	for i := range userItems {
		// Skip if item doesn't exist (pode acontecer se o item não foi criado)
		if userItems[i].ItemID == 0 {
			continue
		}

		// Create user item
		if err := seeders.CreateUserItemIfNotExists(db, &userItems[i]); err != nil {
			log.Fatalf("❌ Failed to create user item: %v", err)
		}
		createdCount++
	}

	fmt.Printf("  ✓ Created user items for %d items\n", createdCount)
}

// printSummary exibe um resumo da operação de seed
func printSummary(totalTags, totalItems int) {
	fmt.Println("═══════════════════════════════════════════")
	fmt.Println("✅ Database seeding completed successfully!")
	fmt.Println("═══════════════════════════════════════════")
	fmt.Println()
	fmt.Println("📊 Summary:")
	fmt.Printf("  • Users:         1\n")
	fmt.Printf("  • Tags:          %d\n", totalTags)
	fmt.Printf("  • Catalog Items: %d\n", totalItems)
	fmt.Printf("  • User Lists:    Populated\n")
	fmt.Println()
	fmt.Println("🎯 You can now:")
	fmt.Println("  1. Login with: demo@geekery.com / password123")
	fmt.Println("  2. Explore the API at: http://localhost:8080/swagger/index.html")
	fmt.Println("  3. Check your personal list at: GET /api/my-list")
	fmt.Println()
}
