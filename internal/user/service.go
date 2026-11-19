package user

import (
	"context"
	"errors"
	"strings"

	"github.com/Sujon37/school-management-system/internal/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"modernc.org/sqlite"
)

var (
	ErrInvalidRole  = errors.New("invalid user role")
	ErrInvalidInput = errors.New("email and password are required")
	ErrEmailExists  = errors.New("email already exists")
)

type CreateUserInput struct {
	Email    string
	Password string
	Role     string
}

type UserService struct {
	userRepo *auth.UserRepository
}

func NewUserService(userRepo *auth.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*auth.User, error) {

	input.Email = strings.TrimSpace(input.Email)
	if input.Email == "" || input.Password == "" {
		return nil, ErrInvalidInput
	}

	if input.Role != "teacher" && input.Role != "student" {
		return nil, ErrInvalidRole
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &auth.User{
		ID:           uuid.NewString(),
		Email:        input.Email,
		PasswordHash: string(hashPassword),
		Role:         input.Role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		var sqliteErr *sqlite.Error

		if errors.As(err, &sqliteErr) && sqliteErr.Code() == 2067 {
			return nil, ErrEmailExists
		}
		return nil, err
	}

	return user, nil
}
