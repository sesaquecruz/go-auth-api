package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sesaquecruz/go-auth-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func Test_FindUserUseCase_NewFindUserUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	findUserUseCase := NewFindUserUseCase(userRepository)
	assert.NotNil(t, findUserUseCase)
	assert.Equal(t, userRepository, findUserUseCase.UserRepository)
}

func Test_FindUserUseCase_Execute_WhenUserExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	findUserUseCase := NewFindUserUseCase(userRepository)

	ctx := context.Background()
	userId := uuid.New()
	user := &entity.User{ID: userId, Email: "user@mail.com"}

	userRepository.EXPECT().FindById(ctx, userId).Return(user, nil).Times(1)

	input := FindUserUseCaseInputDTO{ID: userId.String()}

	output, err := findUserUseCase.Execute(ctx, input)
	assert.Nil(t, err)
	assert.Equal(t, user.Email, output.Email)
}

func Test_FindUserUseCase_Execute_WhenUserNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	findUserUseCase := NewFindUserUseCase(userRepository)

	ctx := context.Background()
	userId := uuid.New()

	userRepository.EXPECT().FindById(ctx, userId).Return(nil, sql.ErrNoRows).Times(1)

	input := FindUserUseCaseInputDTO{ID: userId.String()}

	output, err := findUserUseCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrFindUserUserNotExists)
}

func Test_FindUserUseCase_Execute_WhenUserIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	findUserUseCase := NewFindUserUseCase(userRepository)

	ctx := context.Background()
	userId := uuid.New()

	userRepository.EXPECT().FindById(ctx, userId).Return(nil, sql.ErrNoRows).Times(0)

	input := FindUserUseCaseInputDTO{ID: "aseqrfdf"}

	output, err := findUserUseCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrFindUserInvalidData)
}
