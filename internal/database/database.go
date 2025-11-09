package database

import (
	"fmt"
	"log"

	"github.com/rafaelc-rb/geekery-api/internal/config"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB estabelece conexão com o PostgreSQL e executa as migrações
func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	// Configurar logger do GORM
	logLevel := logger.Silent
	if cfg.Environment == "development" {
		logLevel = logger.Info
	}

	// Conectar ao banco de dados
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✓ Database connection established successfully")

	// Executar auto-migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	DB = db
	return db, nil
}

// runMigrations executa as migrações automáticas do GORM
func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Tag{},
		&models.Item{},     // Catálogo global (sem user_id)
		&models.UserItem{}, // Lista pessoal dos usuários
		// Dados específicos por tipo de mídia
		&models.AnimeData{},
		&models.MovieData{},
		&models.GameData{},
		&models.MusicData{},
		&models.BookData{},
		&models.SeriesData{},
	)
	if err != nil {
		return err
	}

	log.Println("✓ Database migrations completed successfully")
	return nil
}

// GetDB retorna a instância do banco de dados
func GetDB() *gorm.DB {
	return DB
}
