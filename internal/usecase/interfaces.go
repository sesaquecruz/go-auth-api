package usecase

import "context"

type CreateUserUseCaseInterface interface {
	Execute(ctx context.Context, input CreateUserUseCaseInputDTO) error
}
