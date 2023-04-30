package usecase

import (
	"context"
	"errors"

	"github.com/sesaquecruz/go-auth-api/internal/entity"
)

var (
	ErrCreateUserInvalidData      = errors.New("invalid data")
	ErrCreateUserEmailAlreadyUsed = errors.New("email already used")
	ErrCreateUserInternalError    = errors.New("internal error")
)

type CreateUserUseCaseInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserUseCase struct {
	UserFactory    entity.UserFactoryInterface
	UserRepository entity.UserRepositoryInterface
}

func NewCreateUserUseCase(uf entity.UserFactoryInterface, ur entity.UserRepositoryInterface) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserFactory:    uf,
		UserRepository: ur,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserUseCaseInputDTO) error {
	user, err := uc.UserFactory.NewUser(input.Email, input.Password)
	if err != nil {
		return ErrCreateUserInvalidData
	}

	_, err = uc.UserRepository.FindByEmail(ctx, input.Email)
	if err == nil {
		return ErrCreateUserEmailAlreadyUsed
	}

	err = uc.UserRepository.Save(ctx, *user)
	if err != nil {
		return ErrCreateUserInternalError
	}

	return nil
}
