package dto

import (
	"fmt"
	"runtime/debug"
)

// Error codes
const (
	ErrCodeValidation         = "INVALID_INPUT"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeDuplicate          = "DUPLICATE_ENTRY"
	ErrCodeInternal           = "INTERNAL_ERROR"
	ErrCodeInvalidID          = "INVALID_ID"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeUserExists         = "USER_EXISTS"
	ErrCodeForbidden          = "FORBIDDEN"
)

// NewErrorResponse cria uma resposta de erro padronizada
func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: message,
		Code:  code,
	}
}

// NewValidationError cria uma resposta de erro de validação com detalhes
func NewValidationError(message string, details map[string]interface{}) ErrorResponse {
	return ErrorResponse{
		Error:   message,
		Code:    ErrCodeValidation,
		Details: details,
	}
}

// NewInternalError cria uma resposta de erro interno, opcionalmente incluindo stack trace
func NewInternalError(err error, includeStack bool) ErrorResponse {
	response := ErrorResponse{
		Error: "Internal server error",
		Code:  ErrCodeInternal,
	}

	if includeStack && err != nil {
		response.Details = map[string]interface{}{
			"error":      err.Error(),
			"stacktrace": string(debug.Stack()),
		}
	}

	return response
}

// NewNotFoundError cria uma resposta de erro para recursos não encontrados
func NewNotFoundError(resource string) ErrorResponse {
	return ErrorResponse{
		Error: fmt.Sprintf("%s not found", resource),
		Code:  ErrCodeNotFound,
	}
}

// NewDuplicateError cria uma resposta de erro para entradas duplicadas
func NewDuplicateError(resource string) ErrorResponse {
	return ErrorResponse{
		Error: fmt.Sprintf("%s already exists", resource),
		Code:  ErrCodeDuplicate,
	}
}
