package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rafaelc-rb/geekery-api/internal/config"
	"github.com/rafaelc-rb/geekery-api/internal/dto"
)

var validate = validator.New()

// validateAndBind valida e faz bind do JSON para o objeto fornecido
func validateAndBind(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}

	if err := validate.Struct(obj); err != nil {
		return err
	}

	return nil
}

// validateID valida e parseia ID de parâmetro da URL
func validateID(c *gin.Context, param string) (uint, error) {
	idStr := c.Param(param)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}
	return uint(id), nil
}

// respondError envia uma resposta de erro padronizada
func respondError(c *gin.Context, status int, code, message string) {
	response := dto.NewErrorResponse(code, message)
	c.JSON(status, response)
}

// respondValidationError envia uma resposta de erro de validação com detalhes
func respondValidationError(c *gin.Context, err error) {
	var details map[string]interface{}

	// Tentar extrair detalhes de erros de validação do validator
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		details = make(map[string]interface{})
		for _, fieldError := range validationErrors {
			// Usar nome do campo JSON ao invés do nome do campo Go
			jsonFieldName := getJSONFieldName(fieldError)
			details[jsonFieldName] = fieldError.Tag()
		}
	}

	response := dto.NewValidationError(err.Error(), details)
	c.JSON(http.StatusBadRequest, response)
}

// getJSONFieldName extrai o nome do campo JSON de um erro de validação
func getJSONFieldName(fieldError validator.FieldError) string {
	// Converter o nome do campo Go para snake_case/camelCase
	// Exemplo: "Username" -> "username", "ItemID" -> "item_id"
	fieldName := fieldError.Field()

	// Converter para lowercase (funciona para a maioria dos campos)
	return strings.ToLower(fieldName[:1]) + fieldName[1:]
}

// respondInternalError envia uma resposta de erro interno
func respondInternalError(c *gin.Context, err error) {
	isDevelopment := config.AppConfig != nil && config.AppConfig.Environment == "development"
	response := dto.NewInternalError(err, isDevelopment)
	c.JSON(http.StatusInternalServerError, response)
}

// respondSuccess envia uma resposta de sucesso padronizada
func respondSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// respondNotFound envia uma resposta de recurso não encontrado
func respondNotFound(c *gin.Context, resource string) {
	response := dto.NewNotFoundError(resource)
	c.JSON(http.StatusNotFound, response)
}

// respondDuplicate envia uma resposta de entrada duplicada
func respondDuplicate(c *gin.Context, resource string) {
	response := dto.NewDuplicateError(resource)
	c.JSON(http.StatusConflict, response)
}
