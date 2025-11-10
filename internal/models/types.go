package models

// MediaType - Enum para tipos de mídia
type MediaType string

const (
	MediaTypeAnime  MediaType = "anime"
	MediaTypeMovie  MediaType = "movie"
	MediaTypeSeries MediaType = "series"
	MediaTypeGame   MediaType = "game"
	MediaTypeComic  MediaType = "comic"
	MediaTypeNovel  MediaType = "novel"
	MediaTypeBook   MediaType = "book"
)

// ValidMediaTypes lista todos os tipos de mídia válidos
var ValidMediaTypes = []MediaType{
	MediaTypeAnime,
	MediaTypeMovie,
	MediaTypeSeries,
	MediaTypeGame,
	MediaTypeComic,
	MediaTypeNovel,
	MediaTypeBook,
}

// IsValid verifica se o tipo de mídia é válido
func (m MediaType) IsValid() bool {
	for _, valid := range ValidMediaTypes {
		if m == valid {
			return true
		}
	}
	return false
}

// String retorna a representação em string do MediaType
func (m MediaType) String() string {
	return string(m)
}

// MediaStatus - Enum para status do item na lista do usuário
// Usa termos genéricos que funcionam para todos os tipos de mídia
// O frontend deve adaptar as labels conforme o tipo (ex: "Assistindo" para anime, "Jogando" para game)
type MediaStatus string

const (
	StatusPlanned    MediaStatus = "planned"     // Quero assistir/jogar/ler/ouvir
	StatusInProgress MediaStatus = "in_progress" // Assistindo/jogando/lendo/ouvindo
	StatusCompleted  MediaStatus = "completed"   // Completado/assistido/jogado/lido
	StatusPaused     MediaStatus = "paused"      // Pausado/em espera
	StatusDropped    MediaStatus = "dropped"     // Abandonado
)

// ValidStatuses lista todos os status válidos
var ValidStatuses = []MediaStatus{
	StatusPlanned,
	StatusInProgress,
	StatusCompleted,
	StatusPaused,
	StatusDropped,
}

// IsValid verifica se o status é válido
func (s MediaStatus) IsValid() bool {
	for _, valid := range ValidStatuses {
		if s == valid {
			return true
		}
	}
	return false
}

// String retorna a representação em string do MediaStatus
func (s MediaStatus) String() string {
	return string(s)
}

// ProgressType define os tipos de progresso possíveis para diferentes mídias
type ProgressType string

const (
	ProgressTypeEpisodic ProgressType = "episodic" // Series, Anime - season/episode
	ProgressTypeReading  ProgressType = "reading"  // Books, Manga, Light Novels - chapter/volume/page (todos opcionais)
	ProgressTypeTime     ProgressType = "time"     // Movies - minutes_watched
	ProgressTypePercent  ProgressType = "percent"  // Games - percent/hours
	ProgressTypeBoolean  ProgressType = "boolean"  // Music - listened/play_count
)

// ValidProgressTypes lista todos os tipos de progresso válidos
var ValidProgressTypes = []ProgressType{
	ProgressTypeEpisodic,
	ProgressTypeReading,
	ProgressTypeTime,
	ProgressTypePercent,
	ProgressTypeBoolean,
}

// IsValid verifica se o tipo de progresso é válido
func (pt ProgressType) IsValid() bool {
	for _, valid := range ValidProgressTypes {
		if pt == valid {
			return true
		}
	}
	return false
}

// String retorna a representação em string do ProgressType
func (pt ProgressType) String() string {
	return string(pt)
}

// GetDefaultProgressType retorna o tipo de progresso padrão para cada tipo de mídia
func GetDefaultProgressType(mediaType MediaType) ProgressType {
	switch mediaType {
	case MediaTypeAnime, MediaTypeSeries:
		return ProgressTypeEpisodic
	case MediaTypeComic, MediaTypeNovel, MediaTypeBook:
		return ProgressTypeReading
	case MediaTypeMovie:
		return ProgressTypeTime
	case MediaTypeGame:
		return ProgressTypePercent
	default:
		return ProgressTypeBoolean
	}
}
