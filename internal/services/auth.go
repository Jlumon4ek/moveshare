package services

import (
	"errors"
	"moveshare/internal/models"
	"moveshare/internal/repository"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidInput = errors.New("invalid input data")
	ErrInvalidCreds = errors.New("invalid credentials")
)

type AuthService interface {
	CreateUser(req models.SignUpRequest) (*models.User, error)
	Authenticate(req models.LoginRequest) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) CreateUser(req models.SignUpRequest) (*models.User, error) {
	if err := s.validateSignUpRequest(req); err != nil {
		return nil, err
	}

	exists, err := s.userRepo.UserExists(req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	return s.userRepo.CreateUser(user)
}

func (s *authService) Authenticate(req models.LoginRequest) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCreds
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, ErrInvalidCreds
	}
	return user, nil
}

func (s *authService) validateSignUpRequest(req models.SignUpRequest) error {
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return ErrInvalidInput
	}
	if !strings.Contains(req.Email, "@") {
		return ErrInvalidInput
	}
	if len(req.Password) < 6 {
		return ErrInvalidInput
	}
	return nil
}
