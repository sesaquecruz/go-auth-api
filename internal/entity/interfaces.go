package entity

import (
	"context"

	"github.com/google/uuid"
)

type UserFactoryInterface interface {
	NewUser(email string, password string) (*User, error)
}

type UserRepositoryInterface interface {
	Save(ctx context.Context, user User) error
	FindById(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
