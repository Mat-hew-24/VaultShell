package service

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"syncpad/services/auth/internal/model"
	"syncpad/services/auth/internal/repository"
)

var (
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
	ErrEmptyField       = errors.New("email and password are required")
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(email, password string) error {
	// Validate input
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return ErrEmptyField
	}

	// Validate email format
	if !isValidEmail(email) {
		return ErrInvalidEmail
	}

	// Validate password strength
	if len(password) < 6 {
		return ErrPasswordTooShort
	}

	// Check if user already exists
	existingUser, _ := s.repo.FindByEmail(email)
	if existingUser != nil {
		return ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &model.User{
		ID:           uuid.NewString(),
		Email:        strings.TrimSpace(strings.ToLower(email)),
		PasswordHash: string(hash),
	}

	return s.repo.Create(user)
}

// isValidEmail validates email format using regex
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func (s *AuthService) Login(email, password string) (*model.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	//USING BCRYPT TO DO ONE WAY HASH
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
