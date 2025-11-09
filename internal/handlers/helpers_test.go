package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/dto"
)

func TestRespondValidationError_JSONFieldNames(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Criar um handler de teste que valida uma struct com campos comuns
	handler := func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Username string `json:"username" binding:"required,min=3"`
			Password string `json:"password" binding:"required,min=8"`
		}

		if err := validateAndBind(c, &req); err != nil {
			respondValidationError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}

	router := gin.New()
	router.POST("/test", handler)

	// Testar com dados inválidos (todos os campos vazios)
	payload := map[string]string{}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	// Verificar o formato da resposta de erro
	var errorResponse dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v. Body: %s", err, w.Body.String())
	}

	// Verificar que tem detalhes
	if errorResponse.Details == nil {
		t.Fatal("Expected details to be present in error response")
	}

	// Verificar que os nomes dos campos estão em camelCase (json names)
	// Deveria ter "email", "username", "password" - NÃO "Email", "Username", "Password"
	expectedFields := []string{"email", "username", "password"}
	for _, field := range expectedFields {
		if _, exists := errorResponse.Details[field]; !exists {
			t.Errorf("Expected field '%s' in error details, got: %v", field, errorResponse.Details)
		}
	}

	// Verificar que NÃO tem campos em PascalCase
	unexpectedFields := []string{"Email", "Username", "Password"}
	for _, field := range unexpectedFields {
		if _, exists := errorResponse.Details[field]; exists {
			t.Errorf("Did not expect field '%s' in error details (should be lowercase), got: %v", field, errorResponse.Details)
		}
	}
}
