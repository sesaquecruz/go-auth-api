package usecase

import "context"

type CreateUserUseCaseInterface interface {
	Execute(ctx context.Context, input CreateUserUseCaseInputDTO) error
}

type AuthUserUseCaseInterface interface {
	Execute(ctx context.Context, input AuthUserUseCaseInputDTO) (*AuthUserUseCaseOutputDTO, error)
}
