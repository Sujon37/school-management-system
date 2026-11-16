package auth

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInactiveAccount    = errors.New("account is inactive")
	ErrUnauthorizedRole   = errors.New("unauthorized role")
)

type AuthService struct {
	userRepo *UserRepository
}

func NewAuthService(userRepo *UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*User, error) {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !user.IsActive {
		return nil, ErrInactiveAccount
	}

	if user.Role != "admin" {
		return nil, ErrUnauthorizedRole
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}