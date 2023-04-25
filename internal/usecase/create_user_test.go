package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/sesaquecruz/go-auth-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func Test_CreateUserUseCase_NewCreateUserUseCase(t *testing.T) {
	userFactory := &UserFactoryMock{}
	userRepository := &UserRepositoryMock{}

	createUserUseCase := NewCreateUserUseCase(userFactory, userRepository)
	assert.NotNil(t, createUserUseCase)
	assert.Equal(t, userFactory, createUserUseCase.UserFactory)
	assert.Equal(t, userRepository, createUserUseCase.UserRepository)
}

func Test_CreateUserUseCase_Execute_WhenUserIsValid(t *testing.T) {
	user := entity.User{ID: uuid.New(), Email: "user@mail.com", Password: "12345"}
	ctx := context.Background()

	userFactory := &UserFactoryMock{}
	userRepository := &UserRepositoryMock{}

	userFactory.On("NewUser", user.Email, user.Password).Return(&user, nil)
	userRepository.On("FindByEmail", ctx, user.Email).Return(&user, sql.ErrNoRows)
	userRepository.On("Save", ctx, user).Return(nil)

	input := CreateUserUseCaseInputDTO{Email: user.Email, Password: user.Password}
	createUserUseCase := CreateUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	err := createUserUseCase.Execute(ctx, input)
	assert.Nil(t, err)

	userFactory.AssertExpectations(t)
	userFactory.AssertNumberOfCalls(t, "NewUser", 1)
	userFactory.AssertCalled(t, "NewUser", user.Email, user.Password)

	userRepository.AssertExpectations(t)
	userRepository.AssertNumberOfCalls(t, "FindByEmail", 1)
	userRepository.AssertNumberOfCalls(t, "Save", 1)
	userRepository.AssertCalled(t, "FindByEmail", ctx, user.Email)
	userRepository.AssertCalled(t, "Save", ctx, user)
}

func Test_CreateUserUseCase_Execute_WhenUserAlreadyExists(t *testing.T) {
	user := entity.User{ID: uuid.New(), Email: "user@mail.com", Password: "12345"}
	ctx := context.Background()

	userFactory := &UserFactoryMock{}
	userRepository := &UserRepositoryMock{}

	userFactory.On("NewUser", user.Email, user.Password).Return(&user, nil)
	userRepository.On("FindByEmail", ctx, user.Email).Return(&user, nil)
	userRepository.On("Save", ctx, user).Return(nil)

	input := CreateUserUseCaseInputDTO{Email: user.Email, Password: user.Password}
	createUserUseCase := CreateUserUseCase{UserFactory: userFactory, UserRepository: userRepository}

	err := createUserUseCase.Execute(ctx, input)
	assert.ErrorIs(t, err, ErrCreateUserEmailAlreadyUsed)

	userFactory.AssertExpectations(t)
	userFactory.AssertNumberOfCalls(t, "NewUser", 1)
	userFactory.AssertCalled(t, "NewUser", user.Email, user.Password)

	userRepository.AssertNumberOfCalls(t, "FindByEmail", 1)
	userRepository.AssertNumberOfCalls(t, "Save", 0)
	userRepository.AssertCalled(t, "FindByEmail", ctx, user.Email)
}
