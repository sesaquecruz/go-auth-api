package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/sesaquecruz/go-auth-api/internal/entity"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AuthUserUseCase_NewAuthUserUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)

	authUserUseCase := NewAuthUserUseCase(userFactory, userRepository)
	assert.NotNil(t, authUserUseCase)
	assert.Equal(t, userFactory, authUserUseCase.UserFactory)
	assert.Equal(t, userRepository, authUserUseCase.UserRepository)
}

func Test_AuthUserUseCase_Execute_WhenUserIsValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)

	email := "user@mail.com"
	password := "12345"

	ctx := context.Background()
	user, err := entity.NewUserFactory().NewUser(email, password)
	require.Nil(t, err)

	input := AuthUserUseCaseInputDTO{Email: email, Password: password}
	authUserUseCase := AuthUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	userFactory.EXPECT().NewUser(email, password).Return(user, nil).Times(1)
	userRepository.EXPECT().FindByEmail(ctx, email).Return(user, nil).Times(1)

	output, err := authUserUseCase.Execute(ctx, input)
	assert.Nil(t, err)
	assert.Equal(t, output.ID, user.ID.String())
}

func Test_AuthUserUseCase_Execute_WhenUserIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)

	email := "user@mail.com"
	password := "12345"

	ctx := context.Background()

	input := AuthUserUseCaseInputDTO{Email: email, Password: password}
	authUserUseCase := AuthUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	userFactory.EXPECT().NewUser(email, password).Return(nil, errors.New("")).Times(1)

	output, err := authUserUseCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrAuthUserUseCaseInvalidData)
}

func Test_AuthUserUseCase_Execute_WhenUserEmailIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)

	email := "user@mail.com"
	password := "12345"

	ctx := context.Background()
	user, err := entity.NewUserFactory().NewUser(email, password)
	require.Nil(t, err)

	input := AuthUserUseCaseInputDTO{Email: email, Password: password}
	authUserUseCase := AuthUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	userFactory.EXPECT().NewUser(email, password).Return(user, nil)
	userRepository.EXPECT().FindByEmail(ctx, email).Return(nil, sql.ErrNoRows)

	output, err := authUserUseCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrAuthUserUseCaseInvalidCredentials)
}

func Test_AuthUserUseCase_Execute_WhenUserPasswordIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFactory := entity.NewMockUserFactoryInterface(ctrl)
	userRepository := entity.NewMockUserRepositoryInterface(ctrl)

	email := "user@mail.com"
	password := "12345"
	fakePassword := "1234"

	ctx := context.Background()
	user, err := entity.NewUserFactory().NewUser(email, password)
	require.Nil(t, err)

	input := AuthUserUseCaseInputDTO{Email: email, Password: fakePassword}
	authUserUseCase := AuthUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	userFactory.EXPECT().NewUser(email, fakePassword).Return(user, nil)
	userRepository.EXPECT().FindByEmail(ctx, email).Return(user, nil)

	output, err := authUserUseCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrAuthUserUseCaseInvalidCredentials)
}
