package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sesaquecruz/go-auth-api/internal/usecase"

	"github.com/go-chi/jwtauth"
)

type UserHandlerInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandlerMessageDTO struct {
	Message string `json:"message"`
}

type UserHandler struct {
	JWTAuth           *jwtauth.JWTAuth
	JWTExpiration     time.Duration
	CreateUserUseCase usecase.CreateUserUseCaseInterface
	AuthUserUseCase   usecase.AuthUserUseCaseInterface
	UpdateUserUseCase usecase.UpdateUserUseCaseInterface
	DeleteUserUseCase usecase.DeleteUserUseCaseInterface
	FindUserUseCase   usecase.FindUserUseCaseInterface
}

func NewUserHandler(
	jwtAuth *jwtauth.JWTAuth,
	jwtExpiration time.Duration,
	createUserUseCase usecase.CreateUserUseCaseInterface,
	authUserUseCase usecase.AuthUserUseCaseInterface,
	updateUserUseCase usecase.UpdateUserUseCaseInterface,
	deleteUserUseCase usecase.DeleteUserUseCaseInterface,
	findUserUseCase usecase.FindUserUseCaseInterface,
) *UserHandler {
	return &UserHandler{
		JWTAuth:           jwtAuth,
		JWTExpiration:     jwtExpiration,
		CreateUserUseCase: createUserUseCase,
		AuthUserUseCase:   authUserUseCase,
		UpdateUserUseCase: updateUserUseCase,
		DeleteUserUseCase: deleteUserUseCase,
		FindUserUseCase:   findUserUseCase,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var data UserHandlerInputDTO
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.CreateUserUseCase.Execute(r.Context(), usecase.CreateUserUseCaseInputDTO{
		Email:    data.Email,
		Password: data.Password},
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if err == usecase.ErrCreateUserInternalError {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(UserHandlerMessageDTO{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) AuthUser(w http.ResponseWriter, r *http.Request) {
	var data UserHandlerInputDTO
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := h.AuthUserUseCase.Execute(r.Context(), usecase.AuthUserUseCaseInputDTO{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if err == usecase.ErrAuthUserUseCaseInternalError {
			w.WriteHeader(http.StatusInternalServerError)
		} else if err == usecase.ErrAuthUserUseCaseInvalidCredentials {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(UserHandlerMessageDTO{Message: err.Error()})
		return
	}

	payload := map[string]interface{}{
		"sub": output.ID,
		"exp": jwtauth.ExpireIn(h.JWTExpiration),
	}

	_, token, err := h.JWTAuth.Encode(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserHandlerMessageDTO{Message: err.Error()})
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	sub, ok := claims["sub"].(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var data UserHandlerInputDTO
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.UpdateUserUseCase.Execute(r.Context(), usecase.UpdateUserUseCaseInputDTO{
		ID:       sub,
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if err == usecase.ErrUpdateUserInternalError {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(UserHandlerMessageDTO{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	sub, ok := claims["sub"].(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := h.DeleteUserUseCase.Execute(r.Context(), usecase.DeleteUserUseCaseInputDTO{
		ID: sub,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if err == usecase.ErrUpdateUserInternalError {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(UserHandlerMessageDTO{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	sub, ok := claims["sub"].(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	output, err := h.FindUserUseCase.Execute(r.Context(), usecase.FindUserUseCaseInputDTO{
		ID: sub,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if err == usecase.ErrFindUserInternalError {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(UserHandlerMessageDTO{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
