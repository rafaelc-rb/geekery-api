package services

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/rafaelc-rb/geekery-api/internal/dto"
	"github.com/rafaelc-rb/geekery-api/internal/models"
)

// ImportItemsFromCSV importa múltiplos items de um arquivo CSV
func (s *ItemService) ImportItemsFromCSV(ctx context.Context, reader io.Reader, mediaType models.MediaType) (*dto.ImportResult, error) {
	// Validar tipo
	if !mediaType.IsValid() {
		return nil, fmt.Errorf("invalid media type: %s", mediaType)
	}

	// Parse CSV
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, errors.New("CSV file is empty or has no data rows")
	}

	// Validar headers
	headers := records[0]
	if err := s.validateHeadersForType(headers, mediaType); err != nil {
		return nil, err
	}

	result := &dto.ImportResult{
		Success:    true,
		MediaType:  string(mediaType),
		TotalLines: len(records) - 1,
		Errors:     []dto.ImportError{},
	}

	// Processar linhas
	err = s.processImportInTransaction(ctx, records[1:], headers, mediaType, result)

	if err != nil {
		return nil, fmt.Errorf("import failed: %w", err)
	}

	return result, nil
}

// processImportInTransaction processa o import (chamado por ImportItemsFromCSV)
func (s *ItemService) processImportInTransaction(ctx context.Context, records [][]string, headers []string, mediaType models.MediaType, result *dto.ImportResult) error {
	for i, record := range records {
		lineNum := i + 2 // Linha real no CSV (1-indexed + header)

		// Parse linha para Item
		item, specificData, tagNames, err := s.parseCSVRecordByType(headers, record, mediaType)
		if err != nil {
			result.Errors = append(result.Errors, dto.ImportError{
				Line:  lineNum,
				Title: getFieldValue(headers, record, "title"),
				Error: err.Error(),
			})
			result.Failed++
			continue
		}

		// Criar item usando o service normal
		if err := s.createItemWithTags(ctx, item, specificData, tagNames); err != nil {
			result.Errors = append(result.Errors, dto.ImportError{
				Line:  lineNum,
				Title: item.Title,
				Error: err.Error(),
			})
			result.Failed++
			continue
		}

		result.Imported++
	}

	return nil
}

// createItemWithTags cria um item com dados específicos e tags
func (s *ItemService) createItemWithTags(ctx context.Context, item *models.Item, specificData interface{}, tagNames []string) error {
	// Validar item
	if err := item.Validate(); err != nil {
		return err
	}

	// Verificar se já existe item com mesmo título e tipo
	existingItems, _, err := s.itemRepo.SearchByTitle(ctx, item.Title, dto.PaginationParams{Page: 1, Limit: 10})
	if err == nil {
		for _, existing := range existingItems {
			// Comparação case-insensitive do título e tipo exato
			if strings.EqualFold(existing.Title, item.Title) && existing.Type == item.Type {
				return fmt.Errorf("item '%s' (type: %s) already exists with ID %d", item.Title, item.Type, existing.ID)
			}
		}
	}

	// Criar item base
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	// Criar dados específicos
	if specificData != nil {
		if err := s.itemRepo.CreateSpecificData(ctx, item.ID, item.Type, specificData); err != nil {
			return fmt.Errorf("failed to create specific data: %w", err)
		}
	}

	// Associar tags (find or create)
	if len(tagNames) > 0 {
		tagIDs := []uint{}
		for _, name := range tagNames {
			name = strings.ToLower(strings.TrimSpace(name))
			if name == "" {
				continue
			}

			// Find or create tag
			tag := &models.Tag{Name: name}
			if err := s.tagRepo.FindOrCreate(ctx, tag); err != nil {
				return fmt.Errorf("failed to find/create tag '%s': %w", name, err)
			}
			tagIDs = append(tagIDs, tag.ID)
		}

		if len(tagIDs) > 0 {
			if err := s.itemRepo.AssociateTags(ctx, item.ID, tagIDs); err != nil {
				return fmt.Errorf("failed to associate tags: %w", err)
			}
		}
	}

	return nil
}

// validateHeadersForType valida se o CSV tem os headers corretos para o tipo
func (s *ItemService) validateHeadersForType(headers []string, mediaType models.MediaType) error {
	required := getRequiredHeadersForType(mediaType)
	headerSet := make(map[string]bool)

	for _, h := range headers {
		headerSet[strings.ToLower(strings.TrimSpace(h))] = true
	}

	for _, req := range required {
		if !headerSet[req] {
			return fmt.Errorf("missing required column '%s' for type %s", req, mediaType)
		}
	}

	return nil
}

// getRequiredHeadersForType retorna headers obrigatórios por tipo
func getRequiredHeadersForType(mediaType models.MediaType) []string {
	common := []string{"title"} // title sempre obrigatório

	switch mediaType {
	case models.MediaTypeAnime:
		return append(common, "episodes", "studio")
	case models.MediaTypeManga:
		return append(common, "chapters", "author")
	case models.MediaTypeMovie:
		return append(common, "runtime", "director")
	case models.MediaTypeSeries:
		return append(common, "seasons", "episodes")
	case models.MediaTypeGame:
		return append(common, "platform", "developer")
	case models.MediaTypeBook, models.MediaTypeLightNovel:
		return append(common, "pages", "author")
	case models.MediaTypeMusic:
		return append(common, "artist", "duration")
	default:
		return common
	}
}

// parseCSVRecordByType faz parse específico por tipo
func (s *ItemService) parseCSVRecordByType(headers, record []string, mediaType models.MediaType) (*models.Item, interface{}, []string, error) {
	// Parse campos comuns
	item := &models.Item{
		Type:        mediaType,
		Title:       getFieldValue(headers, record, "title"),
		Description: getFieldValue(headers, record, "description"),
		CoverURL:    getFieldValue(headers, record, "cover_url"),
	}

	// Validar title
	if item.Title == "" {
		return nil, nil, nil, errors.New("title is required")
	}

	// Parse release_date
	if dateStr := getFieldValue(headers, record, "release_date"); dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			item.ReleaseDate = &date
		} else {
			return nil, nil, nil, fmt.Errorf("invalid release_date format: %s (use YYYY-MM-DD)", dateStr)
		}
	}

	// Parse tags (separadas por |)
	var tagNames []string
	if tagsStr := getFieldValue(headers, record, "tags"); tagsStr != "" {
		for _, tag := range strings.Split(tagsStr, "|") {
			if tag = strings.TrimSpace(tag); tag != "" {
				tagNames = append(tagNames, tag)
			}
		}
	}

	// Parse external_metadata (formato: source:id|source:id)
	if metadataStr := getFieldValue(headers, record, "external_metadata"); metadataStr != "" {
		item.ExternalMetadata = parseExternalMetadata(metadataStr)
	} else {
		item.ExternalMetadata = models.JSONB{}
	}

	// Parse dados específicos por tipo
	var specificData interface{}
	var err error

	switch mediaType {
	case models.MediaTypeAnime:
		specificData, err = parseAnimeData(headers, record)
	case models.MediaTypeManga:
		specificData, err = parseMangaData(headers, record)
	case models.MediaTypeMovie:
		specificData, err = parseMovieData(headers, record)
	case models.MediaTypeSeries:
		specificData, err = parseSeriesData(headers, record)
	case models.MediaTypeGame:
		specificData, err = parseGameData(headers, record)
	case models.MediaTypeBook, models.MediaTypeLightNovel:
		specificData, err = parseBookData(headers, record)
	case models.MediaTypeMusic:
		specificData, err = parseMusicData(headers, record)
	}

	if err != nil {
		return nil, nil, nil, err
	}

	return item, specificData, tagNames, nil
}

// parseExternalMetadata converte "source:id|source:id" em JSONB
func parseExternalMetadata(metadataStr string) models.JSONB {
	metadata := models.JSONB{}

	for _, pair := range strings.Split(metadataStr, "|") {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			source := strings.ToLower(strings.TrimSpace(parts[0]))
			id := strings.TrimSpace(parts[1])
			if source != "" && id != "" {
				metadata[source] = id
			}
		}
	}

	return metadata
}

// Helper para pegar valor do campo por nome
func getFieldValue(headers, record []string, fieldName string) string {
	for i, h := range headers {
		if strings.ToLower(strings.TrimSpace(h)) == fieldName {
			if i < len(record) {
				return strings.TrimSpace(record[i])
			}
		}
	}
	return ""
}

// Parsers específicos por tipo
func parseAnimeData(headers, record []string) (*models.AnimeData, error) {
	data := &models.AnimeData{
		Studio: getFieldValue(headers, record, "studio"),
	}

	if episodesStr := getFieldValue(headers, record, "episodes"); episodesStr != "" {
		episodes, err := strconv.Atoi(episodesStr)
		if err != nil {
			return nil, fmt.Errorf("invalid episodes value: %s", episodesStr)
		}
		if episodes <= 0 {
			return nil, fmt.Errorf("episodes must be positive: %d", episodes)
		}
		data.Episodes = episodes
	} else {
		return nil, errors.New("episodes is required")
	}

	return data, nil
}

func parseMangaData(headers, record []string) (*models.BookData, error) {
	data := &models.BookData{
		Author: getFieldValue(headers, record, "author"),
	}

	if chaptersStr := getFieldValue(headers, record, "chapters"); chaptersStr != "" {
		chapters, err := strconv.Atoi(chaptersStr)
		if err != nil {
			return nil, fmt.Errorf("invalid chapters value: %s", chaptersStr)
		}
		if chapters <= 0 {
			return nil, fmt.Errorf("chapters must be positive: %d", chapters)
		}
		data.Chapters = chapters
	} else {
		return nil, errors.New("chapters is required")
	}

	if volumesStr := getFieldValue(headers, record, "volumes"); volumesStr != "" {
		volumes, err := strconv.Atoi(volumesStr)
		if err != nil {
			return nil, fmt.Errorf("invalid volumes value: %s", volumesStr)
		}
		if volumes < 0 {
			return nil, fmt.Errorf("volumes cannot be negative: %d", volumes)
		}
		data.Volumes = volumes
	}

	return data, nil
}

func parseMovieData(headers, record []string) (*models.MovieData, error) {
	data := &models.MovieData{
		Director: getFieldValue(headers, record, "director"),
	}

	if runtimeStr := getFieldValue(headers, record, "runtime"); runtimeStr != "" {
		runtime, err := strconv.Atoi(runtimeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid runtime value: %s", runtimeStr)
		}
		if runtime <= 0 {
			return nil, fmt.Errorf("runtime must be positive: %d", runtime)
		}
		data.Runtime = runtime
	} else {
		return nil, errors.New("runtime is required")
	}

	return data, nil
}

func parseSeriesData(headers, record []string) (*models.SeriesData, error) {
	data := &models.SeriesData{}
	// Note: network field exists in CSV but not in model, so we skip it

	if seasonsStr := getFieldValue(headers, record, "seasons"); seasonsStr != "" {
		seasons, err := strconv.Atoi(seasonsStr)
		if err != nil {
			return nil, fmt.Errorf("invalid seasons value: %s", seasonsStr)
		}
		if seasons <= 0 {
			return nil, fmt.Errorf("seasons must be positive: %d", seasons)
		}
		data.Seasons = seasons
	} else {
		return nil, errors.New("seasons is required")
	}

	if episodesStr := getFieldValue(headers, record, "episodes"); episodesStr != "" {
		episodes, err := strconv.Atoi(episodesStr)
		if err != nil {
			return nil, fmt.Errorf("invalid episodes value: %s", episodesStr)
		}
		if episodes <= 0 {
			return nil, fmt.Errorf("episodes must be positive: %d", episodes)
		}
		data.Episodes = episodes
	} else {
		return nil, errors.New("episodes is required")
	}

	return data, nil
}

func parseGameData(headers, record []string) (*models.GameData, error) {
	data := &models.GameData{
		Platform:  getFieldValue(headers, record, "platform"),
		Developer: getFieldValue(headers, record, "developer"),
		// Note: publisher field exists in CSV but not in model, so we skip it
	}

	if data.Platform == "" {
		return nil, errors.New("platform is required")
	}

	if data.Developer == "" {
		return nil, errors.New("developer is required")
	}

	return data, nil
}

func parseBookData(headers, record []string) (*models.BookData, error) {
	data := &models.BookData{
		Author: getFieldValue(headers, record, "author"),
		// Note: publisher field exists in CSV but not in model, so we skip it
	}

	if pagesStr := getFieldValue(headers, record, "pages"); pagesStr != "" {
		pages, err := strconv.Atoi(pagesStr)
		if err != nil {
			return nil, fmt.Errorf("invalid pages value: %s", pagesStr)
		}
		if pages <= 0 {
			return nil, fmt.Errorf("pages must be positive: %d", pages)
		}
		data.Pages = pages
	} else {
		return nil, errors.New("pages is required")
	}

	if data.Author == "" {
		return nil, errors.New("author is required")
	}

	return data, nil
}

func parseMusicData(headers, record []string) (*models.MusicData, error) {
	data := &models.MusicData{
		Artist: getFieldValue(headers, record, "artist"),
		Album:  getFieldValue(headers, record, "album"),
	}

	if durationStr := getFieldValue(headers, record, "duration"); durationStr != "" {
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			return nil, fmt.Errorf("invalid duration value: %s", durationStr)
		}
		if duration <= 0 {
			return nil, fmt.Errorf("duration must be positive: %d", duration)
		}
		data.Duration = duration
	} else {
		return nil, errors.New("duration is required")
	}

	if data.Artist == "" {
		return nil, errors.New("artist is required")
	}

	return data, nil
}
