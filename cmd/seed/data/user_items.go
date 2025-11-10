package data

import (
	"github.com/rafaelc-rb/geekery-api/internal/models"
)

// GetUserItems retorna os items da lista pessoal do usuário para seed
// itemIDs deve ser um map com os títulos dos items mapeados para seus IDs
func GetUserItems(userID uint, itemIDs map[string]uint) []models.UserItem {
	return []models.UserItem{
		{
			UserID:          userID,
			ItemID:          itemIDs["Attack on Titan"],
			Status:          models.StatusCompleted,
			Rating:          9.5,
			Favorite:        true,
			Notes:           "Incrível! Um dos melhores animes que já assisti.",
			ProgressType:    models.ProgressTypeEpisodic,
			CompletionCount: 2,
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
			UserID:          userID,
			ItemID:          itemIDs["Death Note"],
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
			UserID:          userID,
			ItemID:          itemIDs["The Matrix"],
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
			UserID:          userID,
			ItemID:          itemIDs["The Legend of Zelda: Breath of the Wild"],
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
			UserID:          userID,
			ItemID:          itemIDs["One Piece"],
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
			UserID:          userID,
			ItemID:          itemIDs["The Witcher 3: Wild Hunt"],
			Status:          models.StatusCompleted,
			Rating:          9.5,
			Favorite:        true,
			Notes:           "Uma obra-prima dos RPGs. História incrível e mundo imersivo.",
			ProgressType:    models.ProgressTypePercent,
			CompletionCount: 1,
			ProgressData: models.JSONB{
				"percent": 100,
				"hours":   120,
				"extras": map[string]interface{}{
					"dlcs_completed": []string{"Hearts of Stone", "Blood and Wine"},
					"achievements":   85,
				},
				"history": []interface{}{
					map[string]interface{}{
						"started_at":  "2024-03-10T00:00:00Z",
						"finished_at": "2024-05-28T00:00:00Z",
					},
				},
			},
		},
		{
			UserID:          userID,
			ItemID:          itemIDs["Berserk"],
			Status:          models.StatusInProgress,
			Rating:          10.0,
			Favorite:        true,
			Notes:           "Masterpiece absoluta. Arte incomparável.",
			ProgressType:    models.ProgressTypeReading,
			CompletionCount: 0,
			ProgressData: models.JSONB{
				"chapter": 364,
				"volume":  40,
				"history": []interface{}{
					map[string]interface{}{
						"started_at":  "2023-06-15T00:00:00Z",
						"finished_at": nil,
					},
				},
			},
		},
		{
			UserID:          userID,
			ItemID:          itemIDs["Solo Leveling"],
			Status:          models.StatusCompleted,
			Rating:          9.0,
			Favorite:        true,
			Notes:           "Manhwa incrível! Arte espetacular.",
			ProgressType:    models.ProgressTypeReading,
			CompletionCount: 1,
			ProgressData: models.JSONB{
				"chapter": 179,
				"history": []interface{}{
					map[string]interface{}{
						"started_at":  "2023-01-10T00:00:00Z",
						"finished_at": "2023-02-25T00:00:00Z",
					},
				},
			},
		},
		{
			UserID:          userID,
			ItemID:          itemIDs["Sword Art Online"],
			Status:          models.StatusCompleted,
			Rating:          8.0,
			Favorite:        false,
			Notes:           "Light novel que popularizou o gênero isekai moderno.",
			ProgressType:    models.ProgressTypeReading,
			CompletionCount: 1,
			ProgressData: models.JSONB{
				"volume": 28,
				"history": []interface{}{
					map[string]interface{}{
						"started_at":  "2022-08-01T00:00:00Z",
						"finished_at": "2023-01-05T00:00:00Z",
					},
				},
			},
		},
		{
			UserID:          userID,
			ItemID:          itemIDs["Stranger Things"],
			Status:          models.StatusPaused,
			Rating:          7.5,
			Favorite:        false,
			Notes:           "Pausei na 3ª temporada. Vou voltar em breve.",
			ProgressType:    models.ProgressTypeEpisodic,
			CompletionCount: 0,
			ProgressData: models.JSONB{
				"season":  3,
				"episode": 5,
				"history": []interface{}{
					map[string]interface{}{
						"started_at":  "2024-08-01T00:00:00Z",
						"finished_at": nil,
					},
				},
			},
		},
	}
}
