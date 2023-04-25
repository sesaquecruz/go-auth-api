package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/sesaquecruz/go-auth-api/internal/entity"
	"github.com/stretchr/testify/mock"
)

type UserFactoryMock struct {
	mock.Mock
}

func (m *UserFactoryMock) NewUser(email string, password string) (*entity.User, error) {
	args := m.Called(email, password)
	return args.Get(0).(*entity.User), args.Error(1)
}

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Save(ctx context.Context, user entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepositoryMock) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, user entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
