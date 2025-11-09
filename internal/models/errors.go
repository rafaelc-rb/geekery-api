package models

import "errors"

// Erros de validação para UserItem
var (
	ErrUserIDRequired         = errors.New("user ID is required")
	ErrItemIDRequired         = errors.New("item ID is required")
	ErrInvalidStatus          = errors.New("invalid status")
	ErrInvalidRating          = errors.New("rating must be between 0 and 10")
	ErrInvalidProgress        = errors.New("progress cannot be negative")
	ErrInvalidProgressType    = errors.New("invalid progress type")
	ErrInvalidCompletionCount = errors.New("completion count cannot be negative")
	ErrDuplicateEntry         = errors.New("item already in user's list")
)

// Erros de validação para Item
var (
	ErrTitleRequired       = errors.New("title is required")
	ErrInvalidMediaType    = errors.New("invalid media type")
	ErrInvalidYear         = errors.New("year must be between 1900 and current year + 5")
)

// Erros de validação para Tag
var (
	ErrDuplicateTag = errors.New("tag with this name already exists")
	ErrTagNameEmpty = errors.New("tag name cannot be empty")
)
