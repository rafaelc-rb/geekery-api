package testutil

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/rafaelc-rb/geekery-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB cria uma conexão de teste com o banco de dados
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Configurar DSN de teste
	dsn := getTestDSN()

	// Conectar ao banco de dados de teste
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Desabilitar logs em testes
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrar schemas
	if err := MigrateTestDB(db); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// TeardownTestDB limpa o banco de dados de teste
func TeardownTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Limpar todas as tabelas
	CleanupTestDB(t, db)

	// Fechar conexão
	sqlDB, err := db.DB()
	if err != nil {
		t.Logf("Warning: failed to get underlying database: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		t.Logf("Warning: failed to close database connection: %v", err)
	}
}

// CleanupTestDB limpa todas as tabelas do banco de dados de teste
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Ordem de limpeza respeitando foreign keys
	tables := []interface{}{
		&models.UserItem{},
		&models.AnimeData{},
		&models.MovieData{},
		&models.SeriesData{},
		&models.GameData{},
		&models.BookData{},
		&models.MusicData{},
		&models.Item{},
		&models.Tag{},
		&models.User{},
	}

	for _, table := range tables {
		if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(table).Error; err != nil {
			t.Logf("Warning: failed to clean table: %v", err)
		}
	}
}

// MigrateTestDB executa as migrações no banco de teste
func MigrateTestDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Tag{},
		&models.Item{},
		&models.AnimeData{},
		&models.MovieData{},
		&models.SeriesData{},
		&models.GameData{},
		&models.BookData{},
		&models.MusicData{},
		&models.UserItem{},
	)
}

// getTestDSN retorna a DSN para banco de teste
func getTestDSN() string {
	host := getEnvOrDefault("TEST_DB_HOST", "localhost")
	port := getEnvOrDefault("TEST_DB_PORT", "5433")
	user := getEnvOrDefault("TEST_DB_USER", "geekery")
	password := getEnvOrDefault("TEST_DB_PASSWORD", "your_secure_password_here")
	dbname := getEnvOrDefault("TEST_DB_NAME", "geekery_test")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)
}

// getEnvOrDefault retorna variável de ambiente ou valor padrão
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SeedTestData popula o banco com dados de teste
func SeedTestData(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Criar usuário de teste
	user := &models.User{
		Email: "test@example.com",
		Name:  "Test User",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to seed test user: %v", err)
	}

	// Criar tags de teste
	tags := []models.Tag{TagFixtures.ActionTag, TagFixtures.AdventureTag}
	if err := db.Create(&tags).Error; err != nil {
		t.Fatalf("Failed to seed test tags: %v", err)
	}

	// Criar items de teste
	items := []models.Item{ItemFixtures.AnimeItem, ItemFixtures.MovieItem}
	if err := db.Create(&items).Error; err != nil {
		t.Fatalf("Failed to seed test items: %v", err)
	}

	log.Printf("Test data seeded successfully: %d users, %d tags, %d items", 1, len(tags), len(items))
}

// AssertTableCount verifica o número de registros em uma tabela
func AssertTableCount(t *testing.T, db *gorm.DB, model interface{}, expected int64) {
	t.Helper()

	var count int64
	if err := db.Model(model).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count table records: %v", err)
	}

	if count != expected {
		t.Errorf("Expected %d records, got %d", expected, count)
	}
}
