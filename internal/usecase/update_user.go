package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sesaquecruz/go-auth-api/internal/entity"
)

var (
	ErrUpdateUserInvalidData      = errors.New("invalid data")
	ErrUpdateUserUserNotExists    = errors.New("user not exists")
	ErrUpdateUserEmailAlreadyUsed = errors.New("email already used")
	ErrUpdateUserInternalError    = errors.New("internal error")
)

type UpdateUserUseCaseInputDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserUseCaseOutputDTO struct {
	ID string `json:"id"`
}

type UpdateUserUseCase struct {
	UserFactory    entity.UserFactoryInterface
	UserRepository entity.UserRepositoryInterface
}

func NewUpdateUserUseCase(uf entity.UserFactoryInterface, ur entity.UserRepositoryInterface) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		UserFactory:    uf,
		UserRepository: ur,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserUseCaseInputDTO) (*UpdateUserUseCaseOutputDTO, error) {
	user, err := uc.UserFactory.GetUser(input.ID, input.Email, input.Password)
	if err != nil {
		return nil, ErrUpdateUserInvalidData
	}

	_, err = uc.UserRepository.FindById(ctx, user.ID)
	if err != nil {
		return nil, ErrUpdateUserUserNotExists
	}

	emailOwner, err := uc.UserRepository.FindByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, ErrUpdateUserInternalError
	}
	if err == nil && user.ID != emailOwner.ID {
		return nil, ErrUpdateUserEmailAlreadyUsed
	}

	err = uc.UserRepository.Update(ctx, *user)
	if err != nil {
		return nil, ErrUpdateUserInternalError
	}

	output := &UpdateUserUseCaseOutputDTO{
		ID: user.ID.String(),
	}

	return output, nil
}
