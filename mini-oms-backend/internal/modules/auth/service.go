package auth

import (
	"errors"
	"mini-oms-backend/internal/config"
	"mini-oms-backend/internal/models"
	"mini-oms-backend/internal/utils"
)

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User        *models.User `json:"user"`
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int          `json:"expires_in"`
}

// Register registers new user
func (s *Service) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if email already exists
	if s.repo.EmailExists(req.Email) {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", // Default role
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(s.cfg, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.cfg.JWTExpiryHours * 3600, // Convert to seconds
	}, nil
}

// Login authenticates user and returns token
func (s *Service) Login(req *LoginRequest) (*AuthResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(s.cfg, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.cfg.JWTExpiryHours * 3600,
	}, nil
}
