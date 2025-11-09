package models

import "time"

// Helper functions para manipular progresso do UserItem

// getInt converte interface{} para int
func getInt(v interface{}) int {
	if i, ok := v.(int); ok {
		return i
	}
	if f, ok := v.(float64); ok {
		return int(f)
	}
	return 0
}

// getFloat converte interface{} para float64
func getFloat(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	if i, ok := v.(int); ok {
		return float64(i)
	}
	return 0
}

// getHistory retorna o array de history do ProgressData
func (ui *UserItem) getHistory() []interface{} {
	if ui.ProgressData == nil {
		return []interface{}{}
	}

	if history, ok := ui.ProgressData["history"].([]interface{}); ok {
		return history
	}

	return []interface{}{}
}

// SetEpisodicProgress atualiza progresso de séries/anime
func (ui *UserItem) SetEpisodicProgress(season, episode int) {
	ui.ProgressType = ProgressTypeEpisodic

	history := ui.getHistory()

	ui.ProgressData = JSONB{
		"season":  season,
		"episode": episode,
		"history": history,
	}
}

// SetReadingProgress atualiza progresso de livros/manga/light novels
// Todos os parâmetros são opcionais (use ponteiros)
func (ui *UserItem) SetReadingProgress(chapter, volume, page *int) {
	ui.ProgressType = ProgressTypeReading

	history := ui.getHistory()

	data := JSONB{
		"history": history,
	}

	// Adiciona apenas os campos fornecidos
	if chapter != nil {
		data["chapter"] = *chapter
	}
	if volume != nil {
		data["volume"] = *volume
	}
	if page != nil {
		data["page"] = *page
	}

	ui.ProgressData = data
}

// SetTimeProgress atualiza progresso de filmes (minutos assistidos)
func (ui *UserItem) SetTimeProgress(minutesWatched int) {
	ui.ProgressType = ProgressTypeTime

	history := ui.getHistory()

	ui.ProgressData = JSONB{
		"minutes_watched": minutesWatched,
		"last_position":   minutesWatched,
		"history":         history,
	}
}

// SetPercentProgress atualiza progresso de games
func (ui *UserItem) SetPercentProgress(percent int, hours int, extras map[string]interface{}) {
	ui.ProgressType = ProgressTypePercent

	history := ui.getHistory()

	ui.ProgressData = JSONB{
		"percent": percent,
		"hours":   hours,
		"extras":  extras,
		"history": history,
	}
}

// SetBooleanProgress atualiza progresso de música (sem history, só contador)
func (ui *UserItem) SetBooleanProgress(listened bool) {
	ui.ProgressType = ProgressTypeBoolean

	playCount := 0
	if ui.ProgressData != nil {
		if pc, ok := ui.ProgressData["play_count"].(float64); ok {
			playCount = int(pc)
		}
	}

	if listened {
		playCount++
	}

	ui.ProgressData = JSONB{
		"listened":       listened,
		"play_count":     playCount,
		"last_played_at": time.Now(),
	}
}

// StartNewView inicia nova visualização/leitura/playthrough
func (ui *UserItem) StartNewView() {
	history := ui.getHistory()

	// Adiciona nova entrada no history
	newEntry := map[string]interface{}{
		"started_at":  time.Now(),
		"finished_at": nil,
	}
	history = append(history, newEntry)

	// Reseta progresso baseado no tipo
	switch ui.ProgressType {
	case ProgressTypeEpisodic:
		ui.ProgressData = JSONB{
			"season":  1,
			"episode": 0,
			"history": history,
		}
	case ProgressTypeReading:
		ui.ProgressData = JSONB{
			"chapter": 0,
			"volume":  1,
			"page":    0,
			"history": history,
		}
	case ProgressTypeTime:
		ui.ProgressData = JSONB{
			"minutes_watched": 0,
			"last_position":   0,
			"history":         history,
		}
	case ProgressTypePercent:
		ui.ProgressData = JSONB{
			"percent": 0,
			"hours":   0,
			"extras":  map[string]interface{}{},
			"history": history,
		}
	case ProgressTypeBoolean:
		// Boolean não usa history, só incrementa contador no CompletionCount
		ui.ProgressData = JSONB{
			"listened":       false,
			"play_count":     ui.CompletionCount,
			"last_played_at": nil,
		}
	}

	ui.Status = StatusInProgress
}

// CompleteCurrentView finaliza a visualização/leitura/playthrough atual
func (ui *UserItem) CompleteCurrentView() {
	// Para tipos com history
	if ui.ProgressType != ProgressTypeBoolean {
		history := ui.getHistory()

		// Atualiza a última entrada (current view) com finished_at
		if len(history) > 0 {
			lastIndex := len(history) - 1
			if entry, ok := history[lastIndex].(map[string]interface{}); ok {
				if entry["finished_at"] == nil {
					entry["finished_at"] = time.Now()
					history[lastIndex] = entry
				}
			}
		}

		// Atualiza history no ProgressData
		if ui.ProgressData == nil {
			ui.ProgressData = JSONB{}
		}
		ui.ProgressData["history"] = history
	}

	// Incrementa contador de conclusões
	ui.CompletionCount++
	ui.Status = StatusCompleted
}

// GetCurrentViewNumber retorna o número da visualização atual
func (ui *UserItem) GetCurrentViewNumber() int {
	history := ui.getHistory()
	return len(history)
}

// GetCurrentViewStartedAt retorna quando a visualização atual começou
func (ui *UserItem) GetCurrentViewStartedAt() *time.Time {
	history := ui.getHistory()
	if len(history) == 0 {
		return nil
	}

	lastEntry := history[len(history)-1]
	if entry, ok := lastEntry.(map[string]interface{}); ok {
		if startedAt, ok := entry["started_at"].(string); ok {
			t, _ := time.Parse(time.RFC3339, startedAt)
			return &t
		}
		// Também tenta como time.Time
		if startedAt, ok := entry["started_at"].(time.Time); ok {
			return &startedAt
		}
	}

	return nil
}

// IsCurrentViewInProgress verifica se há uma visualização em progresso
func (ui *UserItem) IsCurrentViewInProgress() bool {
	history := ui.getHistory()
	if len(history) == 0 {
		return false
	}

	lastEntry := history[len(history)-1]
	if entry, ok := lastEntry.(map[string]interface{}); ok {
		return entry["finished_at"] == nil
	}

	return false
}

// GetProgressPercent calcula porcentagem baseado no Item (precisa ter Item carregado)
func (ui *UserItem) GetProgressPercent() float64 {
	if ui.ProgressData == nil || ui.Item.ID == 0 {
		return 0
	}

	switch ui.ProgressType {
	case ProgressTypeEpisodic:
		episode := getInt(ui.ProgressData["episode"])
		var total int

		if ui.Item.AnimeData != nil {
			total = ui.Item.AnimeData.Episodes
		} else if ui.Item.SeriesData != nil {
			total = ui.Item.SeriesData.Episodes
		}

		if total > 0 {
			return float64(episode) / float64(total) * 100
		}

	case ProgressTypeTime:
		watched := getInt(ui.ProgressData["minutes_watched"])
		var runtime int

		if ui.Item.MovieData != nil {
			runtime = ui.Item.MovieData.Runtime
		}

		if runtime > 0 {
			return float64(watched) / float64(runtime) * 100
		}

	case ProgressTypeReading:
		// Prioridade: chapter > page > volume
		if ui.Item.BookData != nil {
			// 1. Se tem chapter e total de chapters
			if chapter := getInt(ui.ProgressData["chapter"]); chapter > 0 && ui.Item.BookData.Chapters > 0 {
				return float64(chapter) / float64(ui.Item.BookData.Chapters) * 100
			}

			// 2. Se tem page e total de pages
			if page := getInt(ui.ProgressData["page"]); page > 0 && ui.Item.BookData.Pages > 0 {
				return float64(page) / float64(ui.Item.BookData.Pages) * 100
			}

			// 3. Se tem volume e total de volumes
			if volume := getInt(ui.ProgressData["volume"]); volume > 0 && ui.Item.BookData.Volumes > 0 {
				return float64(volume) / float64(ui.Item.BookData.Volumes) * 100
			}
		}

	case ProgressTypePercent:
		return getFloat(ui.ProgressData["percent"])

	case ProgressTypeBoolean:
		if listened, ok := ui.ProgressData["listened"].(bool); ok && listened {
			return 100
		}
		return 0
	}

	return 0
}

// GetAllViews retorna todas as visualizações/leituras completas
func (ui *UserItem) GetAllViews() []map[string]interface{} {
	history := ui.getHistory()
	result := make([]map[string]interface{}, 0)

	for _, v := range history {
		if entry, ok := v.(map[string]interface{}); ok {
			result = append(result, entry)
		}
	}

	return result
}

// IsRewatching verifica se está re-assistindo/re-lendo
func (ui *UserItem) IsRewatching() bool {
	return ui.CompletionCount > 0 && ui.Status == StatusInProgress
}
