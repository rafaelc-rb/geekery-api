package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rafaelc-rb/geekery-api/internal/config"
	"github.com/rafaelc-rb/geekery-api/internal/database"
	"github.com/rafaelc-rb/geekery-api/internal/models"
)

// parseDate helper para converter string em *time.Time
func parseDate(dateStr string) *time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	return &t
}

func main() {
	// Carregar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Conectar ao banco de dados
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("🌱 Seeding database...")

	// Criar usuário mock
	user := models.User{
		Email: "user@geekery.com",
		Name:  "Geekery User",
	}
	if err := db.FirstOrCreate(&user, models.User{Email: user.Email}).Error; err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("✓ User created: %s (ID: %d)\n", user.Email, user.ID)

	// Criar tags
	tags := []models.Tag{
		{Name: "Action"},
		{Name: "Adventure"},
		{Name: "Drama"},
		{Name: "Comedy"},
		{Name: "Fantasy"},
		{Name: "Sci-Fi"},
		{Name: "Shounen"},
		{Name: "Seinen"},
	}

	for i := range tags {
		if err := db.FirstOrCreate(&tags[i], models.Tag{Name: tags[i].Name}).Error; err != nil {
			log.Fatalf("Failed to create tag: %v", err)
		}
		fmt.Printf("✓ Tag created: %s (ID: %d)\n", tags[i].Name, tags[i].ID)
	}

	// Criar items de catálogo com dados específicos
	items := []struct {
		item         models.Item
		specificData interface{}
	}{
		{
			item: models.Item{
				Title:       "Attack on Titan",
				Type:        models.MediaTypeAnime,
				Description: "A história segue Eren Yeager e seus amigos em um mundo onde a humanidade vive dentro de cidades cercadas por enormes muralhas devido aos Titãs gigantes.",
				ReleaseDate: parseDate("2013-04-07"), // 7 de abril de 2013
				CoverURL:    "https://cdn.myanimelist.net/images/anime/10/47347.jpg",
			},
			specificData: &models.AnimeData{
				Episodes: 75,
				Studio:   "MAPPA, Wit Studio",
			},
		},
		{
			item: models.Item{
				Title:       "Death Note",
				Type:        models.MediaTypeAnime,
				Description: "Light Yagami encontra um caderno sobrenatural que permite matar qualquer pessoa cujo nome seja escrito nele.",
				ReleaseDate: parseDate("2006-10-04"), // 4 de outubro de 2006
				CoverURL:    "https://cdn.myanimelist.net/images/anime/9/9453.jpg",
			},
			specificData: &models.AnimeData{
				Episodes: 37,
				Studio:   "Madhouse",
			},
		},
		{
			item: models.Item{
				Title:       "The Matrix",
				Type:        models.MediaTypeMovie,
				Description: "Um hacker descobre que a realidade é uma simulação criada por máquinas inteligentes.",
				ReleaseDate: parseDate("1999-03-31"), // 31 de março de 1999
				CoverURL:    "https://image.tmdb.org/t/p/w500/f89U3ADr1oiB1s9GkdPOEpXUk5H.jpg",
			},
			specificData: &models.MovieData{
				Director: "Lana Wachowski, Lilly Wachowski",
				Runtime:  136,
			},
		},
		{
			item: models.Item{
				Title:       "Breaking Bad",
				Type:        models.MediaTypeSeries,
				Description: "Um professor de química terminal de câncer se torna um fabricante de metanfetamina.",
				ReleaseDate: parseDate("2008-01-20"), // 20 de janeiro de 2008
				CoverURL:    "https://image.tmdb.org/t/p/w500/ggFHVNu6YYI5L9pCfOacjizRGt.jpg",
			},
			specificData: &models.SeriesData{
				Seasons:  5,
				Episodes: 62,
			},
		},
		{
			item: models.Item{
				Title:       "The Legend of Zelda: Breath of the Wild",
				Type:        models.MediaTypeGame,
				Description: "Link acorda de um sono de 100 anos para salvar Hyrule da Calamidade Ganon.",
				ReleaseDate: parseDate("2017-03-03"), // 3 de março de 2017
				CoverURL:    "https://assets.nintendo.com/image/upload/f_auto/q_auto/dpr_2.0/c_scale,w_500/ncom/en_US/games/switch/t/the-legend-of-zelda-breath-of-the-wild-switch/hero",
			},
			specificData: &models.GameData{
				Platform:        "Nintendo Switch",
				Developer:       "Nintendo EPD",
				AveragePlaytime: 100, // 100 horas médias
			},
		},
		{
			item: models.Item{
				Title:       "One Piece",
				Type:        models.MediaTypeManga,
				Description: "As aventuras de Monkey D. Luffy e sua tripulação em busca do tesouro One Piece.",
				ReleaseDate: parseDate("1997-07-22"),
				CoverURL:    "https://cdn.myanimelist.net/images/manga/2/253146.jpg",
			},
			specificData: &models.BookData{
				Author:   "Eiichiro Oda",
				Volumes:  109,
				Chapters: 1100,
				Pages:    0, // Manga não conta páginas totais
			},
		},
		{
			item: models.Item{
				Title:       "Dark Side of the Moon",
				Type:        models.MediaTypeMusic,
				Description: "Oitavo álbum de estúdio da banda inglesa Pink Floyd.",
				ReleaseDate: parseDate("1973-03-01"),
				CoverURL:    "https://upload.wikimedia.org/wikipedia/en/3/3b/Dark_Side_of_the_Moon.png",
			},
			specificData: &models.MusicData{
				Artist:   "Pink Floyd",
				Album:    "The Dark Side of the Moon",
				Duration: 2583, // 43 minutos em segundos
				Tracks:   10,
			},
		},
	}

	for i := range items {
		// Check if item already exists by title
		var existing models.Item
		if err := db.Where("title = ?", items[i].item.Title).First(&existing).Error; err == nil {
			fmt.Printf("⚠ Item already exists: %s (ID: %d)\n", existing.Title, existing.ID)
			items[i].item = existing
			continue
		}

		// Create item
		if err := db.Create(&items[i].item).Error; err != nil {
			log.Fatalf("Failed to create item: %v", err)
		}

		// Create specific data
		switch items[i].item.Type {
		case models.MediaTypeAnime:
			animeData := items[i].specificData.(*models.AnimeData)
			animeData.ItemID = items[i].item.ID
			if err := db.Create(animeData).Error; err != nil {
				log.Fatalf("Failed to create anime data: %v", err)
			}
		case models.MediaTypeMovie:
			movieData := items[i].specificData.(*models.MovieData)
			movieData.ItemID = items[i].item.ID
			if err := db.Create(movieData).Error; err != nil {
				log.Fatalf("Failed to create movie data: %v", err)
			}
		case models.MediaTypeGame:
			gameData := items[i].specificData.(*models.GameData)
			gameData.ItemID = items[i].item.ID
			if err := db.Create(gameData).Error; err != nil {
				log.Fatalf("Failed to create game data: %v", err)
			}
		case models.MediaTypeSeries:
			seriesData := items[i].specificData.(*models.SeriesData)
			seriesData.ItemID = items[i].item.ID
			if err := db.Create(seriesData).Error; err != nil {
				log.Fatalf("Failed to create series data: %v", err)
			}
		case models.MediaTypeManga, models.MediaTypeLightNovel, models.MediaTypeBook:
			bookData := items[i].specificData.(*models.BookData)
			bookData.ItemID = items[i].item.ID
			if err := db.Create(bookData).Error; err != nil {
				log.Fatalf("Failed to create book data: %v", err)
			}
		case models.MediaTypeMusic:
			musicData := items[i].specificData.(*models.MusicData)
			musicData.ItemID = items[i].item.ID
			if err := db.Create(musicData).Error; err != nil {
				log.Fatalf("Failed to create music data: %v", err)
			}
		}

		// Associar tags apropriadas
		var itemTags []models.Tag
		switch items[i].item.Type {
		case models.MediaTypeAnime:
			itemTags = []models.Tag{tags[0], tags[1], tags[2], tags[6]} // Action, Adventure, Drama, Shounen
		case models.MediaTypeMovie:
			itemTags = []models.Tag{tags[0], tags[5]} // Action, Sci-Fi
		case models.MediaTypeSeries:
			itemTags = []models.Tag{tags[2], tags[3]} // Drama, Comedy
		case models.MediaTypeGame:
			itemTags = []models.Tag{tags[1], tags[4]} // Adventure, Fantasy
		case models.MediaTypeManga:
			itemTags = []models.Tag{tags[0], tags[1], tags[6]} // Action, Adventure, Shounen
		case models.MediaTypeMusic:
			itemTags = []models.Tag{} // Sem tags específicas
		}

		if err := db.Model(&items[i].item).Association("Tags").Replace(itemTags); err != nil {
			log.Fatalf("Failed to associate tags: %v", err)
		}

		fmt.Printf("✓ Item created: %s (ID: %d, Type: %s)\n", items[i].item.Title, items[i].item.ID, items[i].item.Type)
	}

	// Criar alguns user items (itens na lista do usuário) com sistema de progresso
	userItems := []models.UserItem{
		{
			UserID:          user.ID,
			ItemID:          items[0].item.ID, // Attack on Titan
			Status:          models.StatusCompleted,
			Rating:          9.5,
			Favorite:        true,
			Notes:           "Incrível! Um dos melhores animes que já assisti.",
			ProgressType:    models.ProgressTypeEpisodic,
			CompletionCount: 2, // Reassistiu 2 vezes
			ProgressData: models.JSONB{
				"season":  5,
				"episode": 75,
				"history": []interface{}{
					map[string]interface{}{
						"started_at":  "2023-01-15T10:00:00Z",
						"finished_at": "2023-02-20T22:30:00Z",
					},
					map[string]interface{}{
						"started_at":  "2024-06-10T14:00:00Z",
						"finished_at": "2024-07-05T20:00:00Z",
					},
				},
			},
		},
		{
			UserID:          user.ID,
			ItemID:          items[1].item.ID, // Death Note
			Status:          models.StatusInProgress,
			Rating:          8.5,
			Favorite:        false,
			Notes:           "Muito bom até agora!",
			ProgressType:    models.ProgressTypeEpisodic,
			CompletionCount: 0,
			ProgressData: models.JSONB{
				"season":  1,
				"episode": 20,
				"history": []interface{}{
					map[string]interface{}{
						"started_at":  "2025-11-01T08:00:00Z",
						"finished_at": nil,
					},
				},
			},
		},
		{
			UserID:          user.ID,
			ItemID:          items[2].item.ID, // The Matrix
			Status:          models.StatusCompleted,
			Rating:          10.0,
			Favorite:        true,
			Notes:           "Clássico absoluto do cinema.",
			ProgressType:    models.ProgressTypeTime,
			CompletionCount: 1,
			ProgressData: models.JSONB{
				"minutes_watched": 136,
				"last_position":   136,
				"history": []interface{}{
					map[string]interface{}{
						"started_at":  "2024-12-25T20:00:00Z",
						"finished_at": "2024-12-25T22:16:00Z",
					},
				},
			},
		},
		{
			UserID:          user.ID,
			ItemID:          items[4].item.ID, // Zelda BOTW
			Status:          models.StatusInProgress,
			Rating:          9.0,
			Favorite:        false,
			Notes:           "Jogando pela primeira vez, está incrível!",
			ProgressType:    models.ProgressTypePercent,
			CompletionCount: 0,
			ProgressData: models.JSONB{
				"percent": 65,
				"hours":   75,
				"extras": map[string]interface{}{
					"shrines_completed": 85,
					"korok_seeds":       450,
				},
				"history": []interface{}{
					map[string]interface{}{
						"started_at":  "2025-10-15T00:00:00Z",
						"finished_at": nil,
					},
				},
			},
		},
		{
			UserID:          user.ID,
			ItemID:          items[5].item.ID, // One Piece
			Status:          models.StatusInProgress,
			Rating:          9.8,
			Favorite:        true,
			Notes:           "Melhor mangá de todos os tempos!",
			ProgressType:    models.ProgressTypeReading,
			CompletionCount: 0,
			ProgressData: models.JSONB{
				"chapter": 1060,
				"volume":  105,
				"history": []interface{}{
					map[string]interface{}{
						"started_at":  "2020-01-01T00:00:00Z",
						"finished_at": nil,
					},
				},
			},
		},
		{
			UserID:          user.ID,
			ItemID:          items[6].item.ID, // Dark Side of the Moon
			Status:          models.StatusCompleted,
			Rating:          10.0,
			Favorite:        true,
			Notes:           "Obra-prima atemporal.",
			ProgressType:    models.ProgressTypeBoolean,
			CompletionCount: 25, // Ouviu 25 vezes
			ProgressData: models.JSONB{
				"listened":       true,
				"play_count":     25,
				"last_played_at": "2025-11-08T15:30:00Z",
			},
		},
	}

	for i := range userItems {
		// Check if user item already exists
		var existing models.UserItem
		if err := db.Where("user_id = ? AND item_id = ?", userItems[i].UserID, userItems[i].ItemID).First(&existing).Error; err == nil {
			fmt.Printf("⚠ User item already exists: User %d - Item %d\n", existing.UserID, existing.ItemID)
			continue
		}

		if err := db.Create(&userItems[i]).Error; err != nil {
			log.Fatalf("Failed to create user item: %v", err)
		}
		fmt.Printf("✓ User item created: User %d - Item %d (Status: %s)\n", userItems[i].UserID, userItems[i].ItemID, userItems[i].Status)
	}

	fmt.Println("\n✓ Seeding completed successfully!")
	fmt.Println("\n📊 Summary:")
	fmt.Printf("  - %d users\n", 1)
	fmt.Printf("  - %d tags\n", len(tags))
	fmt.Printf("  - %d catalog items\n", len(items))
	fmt.Printf("  - %d user items\n", len(userItems))

	// Show example of specific data
	if len(items) > 0 {
		fmt.Printf("\nExample item with specific data: %s\n", items[0].item.Title)
		fmt.Printf("  Type: %s\n", items[0].item.Type)
		if items[0].item.Type == models.MediaTypeAnime {
			animeData := items[0].specificData.(*models.AnimeData)
			fmt.Printf("  Episodes: %d\n", animeData.Episodes)
			fmt.Printf("  Studio: %s\n", animeData.Studio)
		}
	}
}
