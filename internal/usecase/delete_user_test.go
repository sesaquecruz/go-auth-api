package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/sesaquecruz/go-auth-api/internal/entity"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteUserUseCase_NewDeleteUserUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	deleteUserUseCase := NewDeleteUserUseCase(userRepository)
	assert.NotNil(t, deleteUserUseCase)
	assert.Equal(t, userRepository, deleteUserUseCase.UserRepository)
}

func Test_DeleteUserUseCase_Execute_WhenUserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	deleteUserUseCase := NewDeleteUserUseCase(userRepository)

	ctx := context.Background()
	userId := uuid.New()
	user := &entity.User{ID: userId}

	userRepository.EXPECT().FindById(ctx, userId).Return(user, nil).Times(1)
	userRepository.EXPECT().Delete(ctx, userId).Return(nil).Times(1)

	input := DeleteUserUseCaseInputDTO{ID: userId.String()}

	err := deleteUserUseCase.Execute(ctx, input)
	assert.Nil(t, err)
}

func Test_DeleteUserUseCase_Execute_WhenUserNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	deleteUserUseCase := NewDeleteUserUseCase(userRepository)

	ctx := context.Background()
	userId := uuid.New()

	userRepository.EXPECT().FindById(ctx, userId).Return(nil, sql.ErrNoRows).Times(1)
	userRepository.EXPECT().Delete(ctx, userId).Return(nil).Times(0)

	input := DeleteUserUseCaseInputDTO{ID: userId.String()}

	err := deleteUserUseCase.Execute(ctx, input)
	assert.ErrorIs(t, err, ErrDeleteUserUserNotExists)
}

func Test_DeleteUserUseCase_Execute_WhenUserIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	deleteUserUseCase := NewDeleteUserUseCase(userRepository)

	ctx := context.Background()
	userId := uuid.New()

	userRepository.EXPECT().FindById(ctx, userId).Return(nil, sql.ErrNoRows).Times(0)
	userRepository.EXPECT().Delete(ctx, userId).Return(nil).Times(0)

	input := DeleteUserUseCaseInputDTO{ID: "fiifsiuofef"}

	err := deleteUserUseCase.Execute(ctx, input)
	assert.ErrorIs(t, err, ErrDeleteUserInvalidData)
}
