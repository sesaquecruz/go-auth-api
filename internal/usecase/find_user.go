package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sesaquecruz/go-auth-api/internal/entity"
)

var (
	ErrFindUserInvalidData   = errors.New("invalid data")
	ErrFindUserUserNotExists = errors.New("user not exists")
	ErrFindUserInternalError = errors.New("internal error")
)

type FindUserUseCaseInputDTO struct {
	ID string `json:"id"`
}

type FindUserUseCaseOutputDTO struct {
	Email string `json:"email"`
}

type FindUserUseCase struct {
	UserRepository entity.UserRepositoryInterface
}

func NewFindUserUseCase(ur entity.UserRepositoryInterface) *FindUserUseCase {
	return &FindUserUseCase{UserRepository: ur}
}

func (uc *FindUserUseCase) Execute(ctx context.Context, input FindUserUseCaseInputDTO) (*FindUserUseCaseOutputDTO, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, ErrFindUserInvalidData
	}

	user, err := uc.UserRepository.FindById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFindUserUserNotExists
		}

		return nil, ErrFindUserInternalError
	}

	output := &FindUserUseCaseOutputDTO{
		Email: user.Email,
	}

	return output, nil
}
