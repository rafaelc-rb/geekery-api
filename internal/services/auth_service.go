package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/rafaelc-rb/geekery-api/internal/auth"
	"github.com/rafaelc-rb/geekery-api/internal/models"
	"github.com/rafaelc-rb/geekery-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid username/email or password")
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters")
)

type AuthService struct {
	userRepo   repositories.UserRepositoryInterface
	jwtManager *auth.JWTManager
}

// NewAuthService cria uma nova instância do serviço de autenticação
func NewAuthService(userRepo repositories.UserRepositoryInterface, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register registra um novo usuário
func (s *AuthService) Register(ctx context.Context, email, username, password, name string) (*models.User, error) {
	// Validar senha
	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}

	// Verificar se email já existe
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	// Verificar se username já existe
	_, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil {
		return nil, ErrUsernameAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}

	// Hash da senha (bcrypt cost 12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Criar usuário
	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
		Name:         name,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login autentica um usuário e retorna um token JWT
// Aceita username ou email como identificador
func (s *AuthService) Login(ctx context.Context, usernameOrEmail, password string) (string, *models.User, error) {
	var user *models.User
	var err error

	// Tentar buscar por username primeiro
	user, err = s.userRepo.GetByUsername(ctx, usernameOrEmail)
	if err != nil {
		// Se não encontrou por username, tentar por email
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err = s.userRepo.GetByEmail(ctx, usernameOrEmail)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return "", nil, ErrInvalidCredentials
				}
				return "", nil, fmt.Errorf("failed to get user: %w", err)
			}
		} else {
			return "", nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	// Verificar senha
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Gerar token JWT
	token, err := s.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}
