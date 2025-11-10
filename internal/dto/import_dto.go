package dto

// ImportResult representa o resultado de uma importação CSV
type ImportResult struct {
	Success    bool          `json:"success"`
	MediaType  string        `json:"media_type"`
	TotalLines int           `json:"total_lines"`
	Imported   int           `json:"imported"`
	Failed     int           `json:"failed"`
	Errors     []ImportError `json:"errors,omitempty"`
}

// ImportError representa um erro específico durante a importação
type ImportError struct {
	Line  int    `json:"line"`
	Title string `json:"title,omitempty"`
	Error string `json:"error"`
}
