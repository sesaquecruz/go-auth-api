package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sesaquecruz/go-auth-api/internal/entity"
)

var (
	ErrDeleteUserInvalidData   = errors.New("invalid data")
	ErrDeleteUserUserNotExists = errors.New("user not exists")
	ErrDeleteUserInternalError = errors.New("internal error")
)

type DeleteUserUseCaseInputDTO struct {
	ID string `json:"id"`
}

type DeleteUserUseCase struct {
	UserRepository entity.UserRepositoryInterface
}

func NewDeleteUserUseCase(ur entity.UserRepositoryInterface) *DeleteUserUseCase {
	return &DeleteUserUseCase{UserRepository: ur}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, input DeleteUserUseCaseInputDTO) error {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return ErrDeleteUserInvalidData
	}

	_, err = uc.UserRepository.FindById(ctx, id)
	if err != nil {
		return ErrDeleteUserUserNotExists
	}

	err = uc.UserRepository.Delete(ctx, id)
	if err != nil {
		return ErrDeleteUserInternalError
	}

	return nil
}
