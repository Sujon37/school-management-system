package seed

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Sujon37/school-management-system/internal/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(ctx context.Context, repo *auth.UserRepository, email string, password string) error {
	_, err := repo.GetByEmail(ctx, email)

	if err == nil {
		return nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	admin := &auth.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "admin",
		IsActive:     true,
	}

	return repo.Create(ctx, admin)
}
