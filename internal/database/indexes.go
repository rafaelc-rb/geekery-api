package database

import (
	"log"

	"gorm.io/gorm"
)

// CreateOptimizedIndexes cria índices compostos para otimizar queries frequentes
// Deve ser chamado após as migrations automáticas
func CreateOptimizedIndexes(db *gorm.DB) error {
	log.Println("Creating optimized composite indexes...")

	indexes := []struct {
		name  string
		query string
	}{
		{
			name: "idx_user_items_user_status",
			query: `CREATE INDEX IF NOT EXISTS idx_user_items_user_status
					ON user_items(user_id, status)
					WHERE deleted_at IS NULL`,
		},
		{
			name: "idx_user_items_user_favorite",
			query: `CREATE INDEX IF NOT EXISTS idx_user_items_user_favorite
					ON user_items(user_id, favorite)
					WHERE deleted_at IS NULL AND favorite = true`,
		},
		{
			name: "idx_items_type_date",
			query: `CREATE INDEX IF NOT EXISTS idx_items_type_date
					ON items(type, release_date)
					WHERE deleted_at IS NULL`,
		},
		{
			name: "idx_items_title_lower",
			query: `CREATE INDEX IF NOT EXISTS idx_items_title_lower
					ON items(LOWER(title))
					WHERE deleted_at IS NULL`,
		},
	}

	for _, idx := range indexes {
		log.Printf("Creating index: %s", idx.name)
		if err := db.Exec(idx.query).Error; err != nil {
			log.Printf("Warning: Failed to create index %s: %v", idx.name, err)
			// Não retornar erro para não quebrar a aplicação
			// Índices são otimizações, não requisitos
		} else {
			log.Printf("✓ Index %s created successfully", idx.name)
		}
	}

	log.Println("✓ Optimized indexes creation completed")
	return nil
}

// CreateFullTextSearchIndex cria índice para full-text search usando pg_trgm
// Requer extensão pg_trgm habilitada no PostgreSQL
// Execute manualmente no PostgreSQL: CREATE EXTENSION IF NOT EXISTS pg_trgm;
func CreateFullTextSearchIndex(db *gorm.DB) error {
	log.Println("Creating full-text search index (requires pg_trgm extension)...")

	// Habilitar extensão pg_trgm
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
		log.Printf("Warning: Could not enable pg_trgm extension: %v", err)
		log.Println("Full-text search will use standard LIKE queries (slower)")
		return nil
	}

	// Criar índice GIN para full-text search
	query := `CREATE INDEX IF NOT EXISTS idx_items_title_trgm
			  ON items USING gin(LOWER(title) gin_trgm_ops)`

	if err := db.Exec(query).Error; err != nil {
		log.Printf("Warning: Failed to create full-text search index: %v", err)
		return nil
	}

	log.Println("✓ Full-text search index created successfully")
	return nil
}
