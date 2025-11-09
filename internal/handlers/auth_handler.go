package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/dto"
	"github.com/rafaelc-rb/geekery-api/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler cria uma nova instância do handler de autenticação
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register registra um novo usuário
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterRequest  true  "Registration data"
// @Success      201      {object}  dto.AuthResponse     "User registered successfully"
// @Failure      400      {object}  map[string]string    "Bad request - validation error"
// @Failure      409      {object}  map[string]string    "Conflict - email or username already exists"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RegisterRequest
	if err := validateAndBind(c, &req); err != nil {
		respondValidationError(c, err)
		return
	}

	// Registrar usuário
	user, err := h.authService.Register(ctx, req.Email, req.Username, req.Password, req.Name)
	if err != nil {
		switch err {
		case services.ErrEmailAlreadyExists:
			respondError(c, http.StatusConflict, dto.ErrCodeUserExists, "email already exists")
		case services.ErrUsernameAlreadyExists:
			respondError(c, http.StatusConflict, dto.ErrCodeUserExists, "username already exists")
		case services.ErrPasswordTooShort:
			respondError(c, http.StatusBadRequest, dto.ErrCodeValidation, err.Error())
		default:
			respondInternalError(c, err)
		}
		return
	}

	// Fazer login automático após registro
	token, _, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		// Se falhar o login, ainda retornar sucesso no registro
		respondSuccess(c, http.StatusCreated, gin.H{
			"message": "user registered successfully, please login",
			"user": dto.UserInfo{
				ID:        user.ID,
				Email:     user.Email,
				Username:  user.Username,
				Name:      user.Name,
				CreatedAt: user.CreatedAt,
			},
		})
		return
	}

	// Retornar token e dados do usuário
	response := dto.AuthResponse{
		Token: token,
		User: dto.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
	}

	respondSuccess(c, http.StatusCreated, response)
}

// Login autentica um usuário
// @Summary      Login
// @Description  Authenticate user with username or email and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest   true  "Login credentials (username or email)"
// @Success      200      {object}  dto.AuthResponse   "Login successful"
// @Failure      400      {object}  map[string]string  "Bad request - validation error"
// @Failure      401      {object}  map[string]string  "Unauthorized - invalid credentials"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.LoginRequest
	if err := validateAndBind(c, &req); err != nil {
		respondValidationError(c, err)
		return
	}

	// Autenticar usuário (username ou email)
	token, user, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			respondError(c, http.StatusUnauthorized, dto.ErrCodeInvalidCredentials, "invalid username/email or password")
			return
		}
		respondInternalError(c, err)
		return
	}

	// Retornar token e dados do usuário
	response := dto.AuthResponse{
		Token: token,
		User: dto.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
	}

	respondSuccess(c, http.StatusOK, response)
}
