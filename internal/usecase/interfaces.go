package usecase

import "context"

type CreateUserUseCaseInterface interface {
	Execute(ctx context.Context, input CreateUserUseCaseInputDTO) error
}

type AuthUserUseCaseInterface interface {
	Execute(ctx context.Context, input AuthUserUseCaseInputDTO) (*AuthUserUseCaseOutputDTO, error)
}

type UpdateUserUseCaseInterface interface {
	Execute(ctx context.Context, input UpdateUserUseCaseInputDTO) (*UpdateUserUseCaseOutputDTO, error)
}

type DeleteUserUseCaseInterface interface {
	Execute(ctx context.Context, input DeleteUserUseCaseInputDTO) error
}
