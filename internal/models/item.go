package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// JSONB é um tipo customizado para armazenar dados JSON no PostgreSQL
type JSONB map[string]interface{}

// Value implementa a interface driver.Valuer para JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implementa a interface sql.Scanner para JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONB value")
	}

	result := make(JSONB)
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// Item representa uma mídia no catálogo global (anime, filme, série, jogo, etc.)
// Este é o catálogo compartilhado - cada item existe apenas UMA vez
// Usuários adicionam items à sua lista pessoal via UserItem
type Item struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Title       string         `json:"title" gorm:"not null;index"`
	Type        MediaType      `json:"type" gorm:"type:varchar(50);not null;index;check:type IN ('anime','movie','series','game','comic','novel','book')"`
	Description string         `json:"description" gorm:"type:text"`
	ReleaseDate *time.Time     `json:"release_date" gorm:"index"` // Data de lançamento/estreia
	CoverURL    string         `json:"cover_url"`
	ExternalMetadata JSONB     `json:"external_metadata" gorm:"type:jsonb"` // Metadados de APIs externas (MAL, IMDb, etc)
	Tags        []Tag          `json:"tags,omitempty" gorm:"many2many:item_tags;"`

	// Dados específicos por tipo (apenas um será não-nil baseado no Type)
	AnimeData  *AnimeData  `json:"anime_data,omitempty" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	MovieData  *MovieData  `json:"movie_data,omitempty" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	GameData   *GameData   `json:"game_data,omitempty" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	BookData   *BookData   `json:"book_data,omitempty" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	SeriesData *SeriesData `json:"series_data,omitempty" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
}

// TableName especifica o nome da tabela no banco de dados
func (Item) TableName() string {
	return "items"
}

// Validate valida os dados do Item
func (i *Item) Validate() error {
	if i.Title == "" {
		return ErrTitleRequired
	}

	if len(i.Title) > 500 {
		return errors.New("title must be at most 500 characters")
	}

	if !i.Type.IsValid() {
		return ErrInvalidMediaType
	}

	// Validação opcional da data de lançamento
	if i.ReleaseDate != nil {
		minDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		maxDate := time.Now().AddDate(5, 0, 0) // Até 5 anos no futuro

		if i.ReleaseDate.Before(minDate) || i.ReleaseDate.After(maxDate) {
			return ErrInvalidYear // Reutilizando erro existente
		}
	}

	return nil
}
