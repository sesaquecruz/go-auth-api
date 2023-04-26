package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sesaquecruz/go-auth-api/internal/usecase"

	"github.com/go-chi/jwtauth"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserHandler_NewUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	createUserUseCase := usecase.NewMockCreateUserUseCaseInterface(ctrl)
	authUserUseCase := usecase.NewMockAuthUserUseCaseInterface(ctrl)

	jwtAuth := jwtauth.New("HS256", []byte("secret"), nil)
	jwtxpiration := time.Duration(300) * time.Second

	userHander := NewUserHandler(
		jwtAuth,
		jwtxpiration,
		createUserUseCase,
		authUserUseCase,
	)
	assert.NotNil(t, userHander)
	assert.Equal(t, jwtAuth, userHander.JWTAuth)
	assert.Equal(t, jwtxpiration, userHander.JWTExpiration)
	assert.Equal(t, createUserUseCase, userHander.CreateUserUseCase)
	assert.Equal(t, authUserUseCase, userHander.AuthUserUseCase)
}

func Test_UserHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	createUserUseCase := usecase.NewMockCreateUserUseCaseInterface(ctrl)

	userHander := UserHandler{
		CreateUserUseCase: createUserUseCase,
	}

	createUserUseCase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	ts := httptest.NewServer(http.HandlerFunc(userHander.CreateUser))
	defer ts.Close()

	data := UserHandlerInputDTO{
		Email:    "user@mail.com",
		Password: "12345",
	}

	contentType := "application/json"
	body, err := json.Marshal(data)
	require.Nil(t, err)

	response, err := http.Post(ts.URL, contentType, bytes.NewReader(body))
	assert.Nil(t, err)
	defer response.Body.Close()

	assert.Equal(t, response.StatusCode, http.StatusCreated)
}

func Test_UserHandler_AuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authUserUseCase := usecase.NewMockAuthUserUseCaseInterface(ctrl)

	userHander := UserHandler{
		JWTAuth:         jwtauth.New("HS256", []byte("secret"), nil),
		JWTExpiration:   time.Duration(300) * time.Second,
		AuthUserUseCase: authUserUseCase,
	}

	output := &usecase.AuthUserUseCaseOutputDTO{ID: uuid.NewString()}
	authUserUseCase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(output, nil).Times(1)

	ts := httptest.NewServer(http.HandlerFunc(userHander.AuthUser))
	defer ts.Close()

	data := UserHandlerInputDTO{
		Email:    "user@mail.com",
		Password: "12345",
	}

	contentType := "application/json"
	body, err := json.Marshal(data)
	require.Nil(t, err)

	response, err := http.Post(ts.URL, contentType, bytes.NewReader(body))
	assert.Nil(t, err)
	defer response.Body.Close()

	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.NotEmpty(t, response.Header.Get("Authorization"))
}
