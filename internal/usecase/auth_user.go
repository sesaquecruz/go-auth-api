package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sesaquecruz/go-auth-api/internal/entity"
)

var (
	ErrAuthUserUseCaseInvalidData        = errors.New("invalid data")
	ErrAuthUserUseCaseInternalError      = errors.New("internal error")
	ErrAuthUserUseCaseInvalidCredentials = errors.New("invalid credentials")
)

type AuthUserUseCaseInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUserUseCaseOutputDTO struct {
	ID string `json:"id"`
}

type AuthUserUseCase struct {
	UserFactory    entity.UserFactoryInterface
	UserRepository entity.UserRepositoryInterface
}

func NewAuthUserUseCase(uf entity.UserFactoryInterface, ur entity.UserRepositoryInterface) *AuthUserUseCase {
	return &AuthUserUseCase{
		UserFactory:    uf,
		UserRepository: ur,
	}
}

func (uc *AuthUserUseCase) Execute(ctx context.Context, input AuthUserUseCaseInputDTO) (*AuthUserUseCaseOutputDTO, error) {
	_, err := uc.UserFactory.NewUser(input.Email, input.Password)
	if err != nil {
		return nil, ErrAuthUserUseCaseInvalidData
	}

	user, err := uc.UserRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAuthUserUseCaseInvalidCredentials
		}
		return nil, ErrAuthUserUseCaseInternalError
	}

	err = user.VerifyPassword(input.Password)
	if err != nil {
		return nil, ErrAuthUserUseCaseInvalidCredentials
	}

	output := &AuthUserUseCaseOutputDTO{
		ID: user.ID.String(),
	}

	return output, nil
}
