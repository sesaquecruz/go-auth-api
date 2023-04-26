package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sesaquecruz/go-auth-api/internal/entity"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CreateUserUseCase_NewCreateUserUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)

	createUserUseCase := NewCreateUserUseCase(userFactory, userRepository)
	assert.NotNil(t, createUserUseCase)
	assert.Equal(t, userFactory, createUserUseCase.UserFactory)
	assert.Equal(t, userRepository, createUserUseCase.UserRepository)
}

func Test_CreateUserUseCase_Execute_WhenUserIsValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)

	user := &entity.User{ID: uuid.New(), Email: "user@mail.com", Password: "12345"}
	ctx := context.Background()

	userFactory.EXPECT().NewUser(user.Email, user.Password).Return(user, nil).Times(1)
	userRepository.EXPECT().FindByEmail(ctx, user.Email).Return(nil, sql.ErrNoRows).Times(1)
	userRepository.EXPECT().Save(ctx, *user).Return(nil).Times(1)

	createUserUseCase := CreateUserUseCase{UserFactory: userFactory, UserRepository: userRepository}
	input := CreateUserUseCaseInputDTO{Email: user.Email, Password: user.Password}

	err := createUserUseCase.Execute(ctx, input)
	assert.Nil(t, err)
}

func Test_CreateUserUseCase_Execute_WhenUserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)

	user := &entity.User{ID: uuid.New(), Email: "user@mail.com", Password: "12345"}
	ctx := context.Background()

	userFactory.EXPECT().NewUser(user.Email, user.Password).Return(user, nil).Times(1)
	userRepository.EXPECT().FindByEmail(ctx, user.Email).Return(user, nil).Times(1)
	userRepository.EXPECT().Save(ctx, *user).Return(nil).Times(0)

	input := CreateUserUseCaseInputDTO{Email: user.Email, Password: user.Password}
	createUserUseCase := CreateUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	err := createUserUseCase.Execute(ctx, input)
	assert.ErrorIs(t, err, ErrCreateUserEmailAlreadyUsed)
}
