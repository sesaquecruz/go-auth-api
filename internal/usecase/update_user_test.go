package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/sesaquecruz/go-auth-api/internal/entity"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateUserUseCase_NewUpdateUserUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	updateUserUseCase := NewUpdateUserUseCase(userFactory, userRepository)
	assert.NotNil(t, updateUserUseCase)
	assert.Equal(t, userFactory, updateUserUseCase.UserFactory)
	assert.Equal(t, userRepository, updateUserUseCase.UserRepository)
}

func Test_UpdateUserUseCase_Execute_WhenUserDataIsValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &entity.User{
		ID:       uuid.New(),
		Email:    "user@mail.com",
		Password: "12345",
	}

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	updateUserUseCase := UpdateUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	ctx := context.Background()
	input := UpdateUserUseCaseInputDTO{
		ID:       user.ID.String(),
		Email:    user.Email,
		Password: user.Password,
	}

	userFactory.EXPECT().GetUser(input.ID, input.Email, input.Password).Return(user, nil).Times(1)
	userRepository.EXPECT().FindById(ctx, user.ID).Return(user, nil).Times(1)
	userRepository.EXPECT().FindByEmail(ctx, input.Email).Return(user, nil).Times(1)
	userRepository.EXPECT().Update(ctx, *user).Return(nil).Times(1)

	output, err := updateUserUseCase.Execute(ctx, input)
	assert.Nil(t, err)
	assert.Equal(t, user.ID.String(), output.ID)
}

func Test_UpdateUserUseCase_Execute_WhenUserDataIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &entity.User{
		ID:       uuid.New(),
		Email:    "user@mailcom",
		Password: "12345",
	}

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	updateUserUseCase := UpdateUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	ctx := context.Background()
	input := UpdateUserUseCaseInputDTO{
		ID:       user.ID.String(),
		Email:    user.Email,
		Password: user.Password,
	}

	userFactory.EXPECT().GetUser(input.ID, input.Email, input.Password).Return(nil, errors.New("")).Times(1)

	output, err := updateUserUseCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrUpdateUserInvalidData)
}

func Test_UpdateUserUseCase_Execute_WhenUserIdIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &entity.User{
		ID:       uuid.New(),
		Email:    "user@mail.com",
		Password: "12345",
	}

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	updateUserUseCase := UpdateUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	ctx := context.Background()
	input := UpdateUserUseCaseInputDTO{
		ID:       user.ID.String(),
		Email:    user.Email,
		Password: user.Password,
	}

	userFactory.EXPECT().GetUser(input.ID, input.Email, input.Password).Return(user, nil).Times(1)
	userRepository.EXPECT().FindById(ctx, user.ID).Return(nil, errors.New("")).Times(1)

	output, err := updateUserUseCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrUpdateUserUserNotExists)
}

func Test_UpdateUserUseCase_Execute_WhenUserEmailIsNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &entity.User{
		ID:       uuid.New(),
		Email:    "user@mail.com",
		Password: "12345",
	}

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	updateUserUseCase := UpdateUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	ctx := context.Background()
	input := UpdateUserUseCaseInputDTO{
		ID:       user.ID.String(),
		Email:    user.Email,
		Password: user.Password,
	}

	userFactory.EXPECT().GetUser(input.ID, input.Email, input.Password).Return(user, nil).Times(1)
	userRepository.EXPECT().FindById(ctx, user.ID).Return(user, nil).Times(1)
	userRepository.EXPECT().FindByEmail(ctx, input.Email).Return(nil, sql.ErrNoRows).Times(1)
	userRepository.EXPECT().Update(ctx, *user).Return(nil).Times(1)

	output, err := updateUserUseCase.Execute(ctx, input)
	assert.Nil(t, err)
	assert.Equal(t, user.ID.String(), output.ID)
}

func Test_UpdateUserUseCase_Execute_WhenUserEmailAlreadyUsed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &entity.User{
		ID:       uuid.New(),
		Email:    "user@mail.com",
		Password: "12345",
	}

	otherUser := &entity.User{
		ID:       uuid.New(),
		Email:    "user@mail.com",
		Password: "12345",
	}

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)
	updateUserUseCase := UpdateUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	ctx := context.Background()
	input := UpdateUserUseCaseInputDTO{
		ID:       user.ID.String(),
		Email:    user.Email,
		Password: user.Password,
	}

	userFactory.EXPECT().GetUser(input.ID, input.Email, input.Password).Return(user, nil).Times(1)
	userRepository.EXPECT().FindById(ctx, user.ID).Return(user, nil).Times(1)
	userRepository.EXPECT().FindByEmail(ctx, input.Email).Return(otherUser, nil).Times(1)

	output, err := updateUserUseCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrUpdateUserEmailAlreadyUsed)
}
