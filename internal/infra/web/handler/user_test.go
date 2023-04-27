package handler

import (
	"bytes"
	"context"
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
	updateUserUsecase := usecase.NewMockUpdateUserUseCaseInterface(ctrl)
	deleteUserUseCase := usecase.NewMockDeleteUserUseCaseInterface(ctrl)

	jwtAuth := jwtauth.New("HS256", []byte("secret"), nil)
	jwtxpiration := time.Duration(300) * time.Second

	userHander := NewUserHandler(
		jwtAuth,
		jwtxpiration,
		createUserUseCase,
		authUserUseCase,
		updateUserUsecase,
		deleteUserUseCase,
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

	assert.Equal(t, http.StatusCreated, response.StatusCode)
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, response.Header.Get("Authorization"))
}

func Test_UserHandler_UpdateUser(t *testing.T) {
	jwtAuth := jwtauth.New("HS256", []byte("secret"), nil)
	payload := map[string]interface{}{
		"sub": uuid.NewString(),
		"exp": jwtauth.ExpireIn(time.Duration(300) * time.Second),
	}
	token, _, err := jwtAuth.Encode(payload)
	require.Nil(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	updateUserUseCase := usecase.NewMockUpdateUserUseCaseInterface(ctrl)
	updateUserUseCase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

	userHandler := UserHandler{UpdateUserUseCase: updateUserUseCase}

	body, err := json.Marshal(UserHandlerInputDTO{Email: "user@mail.com", Password: "12345"})
	require.Nil(t, err)

	ctx := jwtauth.NewContext(context.Background(), token, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/", bytes.NewReader(body))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	userHandler.UpdateUser(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_UserHandler_DelteUser(t *testing.T) {
	jwtAuth := jwtauth.New("HS256", []byte("secret"), nil)
	payload := map[string]interface{}{
		"sub": uuid.NewString(),
		"exp": jwtauth.ExpireIn(time.Duration(300) * time.Second),
	}
	token, _, err := jwtAuth.Encode(payload)
	require.Nil(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deleteUserUseCase := usecase.NewMockDeleteUserUseCaseInterface(ctrl)
	deleteUserUseCase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	userHandler := UserHandler{DeleteUserUseCase: deleteUserUseCase}

	ctx := jwtauth.NewContext(context.Background(), token, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	userHandler.DeleteUser(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
