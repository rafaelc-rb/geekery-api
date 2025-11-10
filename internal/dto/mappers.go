package dto

import (
	"github.com/rafaelc-rb/geekery-api/internal/models"
)

// ItemToDTO converte um Item model para ItemDTO
func ItemToDTO(item *models.Item) *ItemDTO {
	if item == nil {
		return nil
	}

	dto := &ItemDTO{
		ID:               item.ID,
		Title:            item.Title,
		Type:             string(item.Type),
		Description:      item.Description,
		ReleaseDate:      item.ReleaseDate,
		CoverURL:         item.CoverURL,
		ExternalMetadata: item.ExternalMetadata,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}

	// Converter tags
	if len(item.Tags) > 0 {
		dto.Tags = make([]TagDTO, len(item.Tags))
		for i, tag := range item.Tags {
			dto.Tags[i] = *TagToDTO(&tag)
		}
	}

	// Adicionar dados específicos baseado no tipo
	switch item.Type {
	case models.MediaTypeAnime:
		if item.AnimeData != nil {
			dto.SpecificData = &AnimeDataDTO{
				Episodes: item.AnimeData.Episodes,
				Studio:   item.AnimeData.Studio,
			}
		}
	case models.MediaTypeMovie:
		if item.MovieData != nil {
			dto.SpecificData = &MovieDataDTO{
				Director: item.MovieData.Director,
				Runtime:  item.MovieData.Runtime,
			}
		}
	case models.MediaTypeSeries:
		if item.SeriesData != nil {
			dto.SpecificData = &SeriesDataDTO{
				Seasons:  item.SeriesData.Seasons,
				Episodes: item.SeriesData.Episodes,
			}
		}
	case models.MediaTypeGame:
		if item.GameData != nil {
			dto.SpecificData = &GameDataDTO{
				Platform:        item.GameData.Platform,
				Developer:       item.GameData.Developer,
				AveragePlaytime: item.GameData.AveragePlaytime,
			}
		}
	case models.MediaTypeBook, models.MediaTypeComic, models.MediaTypeNovel:
		if item.BookData != nil {
			dto.SpecificData = &BookDataDTO{
				Author:    item.BookData.Author,
				Volumes:   item.BookData.Volumes,
				Chapters:  item.BookData.Chapters,
				Pages:     item.BookData.Pages,
				Format:    item.BookData.Format,
				Publisher: item.BookData.Publisher,
			}
		}
	}

	return dto
}

// ItemsToDTOs converte uma slice de Items para slice de ItemDTOs
func ItemsToDTOs(items []models.Item) []ItemDTO {
	dtos := make([]ItemDTO, len(items))
	for i, item := range items {
		if dto := ItemToDTO(&item); dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}

// UserItemToDTO converte um UserItem model para UserItemDTO
func UserItemToDTO(userItem *models.UserItem) *UserItemDTO {
	if userItem == nil {
		return nil
	}

	dto := &UserItemDTO{
		ID:              userItem.ID,
		UserID:          userItem.UserID,
		ItemID:          userItem.ItemID,
		Status:          string(userItem.Status),
		Rating:          userItem.Rating,
		Favorite:        userItem.Favorite,
		Notes:           userItem.Notes,
		ProgressType:    string(userItem.ProgressType),
		ProgressData:    userItem.ProgressData,
		CompletionCount: userItem.CompletionCount,
		CreatedAt:       userItem.CreatedAt,
		UpdatedAt:       userItem.UpdatedAt,
	}

	// Incluir item completo se disponível
	if userItem.Item.ID != 0 {
		dto.Item = ItemToDTO(&userItem.Item)
	}

	return dto
}

// UserItemsToDTOs converte uma slice de UserItems para slice de UserItemDTOs
func UserItemsToDTOs(userItems []models.UserItem) []UserItemDTO {
	dtos := make([]UserItemDTO, len(userItems))
	for i, userItem := range userItems {
		if dto := UserItemToDTO(&userItem); dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}

// TagToDTO converte um Tag model para TagDTO
func TagToDTO(tag *models.Tag) *TagDTO {
	if tag == nil {
		return nil
	}

	return &TagDTO{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
}

// TagsToDTOs converte uma slice de Tags para slice de TagDTOs
func TagsToDTOs(tags []models.Tag) []TagDTO {
	dtos := make([]TagDTO, len(tags))
	for i, tag := range tags {
		if dto := TagToDTO(&tag); dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}

// StatsToDTO converte um mapa de estatísticas para UserListStatsDTO
func StatsToDTO(stats map[string]int64) *UserListStatsDTO {
	if stats == nil {
		return nil
	}

	return &UserListStatsDTO{
		Total:      stats["total"],
		InProgress: stats["in_progress"],
		Completed:  stats["completed"],
		Planned:    stats["planned"],
		Paused:     stats["paused"],
		Dropped:    stats["dropped"],
		Favorites:  stats["favorites"],
	}
}
