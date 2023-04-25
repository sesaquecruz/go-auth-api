package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sesaquecruz/go-auth-api/internal/usecase"
)

type UserHandlerInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandlerOutputDTO struct {
	Message string `json:"message"`
}

type UserHandler struct {
	CreateUserUseCase usecase.CreateUserUseCaseInterface
}

func NewUserHandler(createUserUseCase usecase.CreateUserUseCaseInterface) *UserHandler {
	return &UserHandler{
		CreateUserUseCase: createUserUseCase,
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

		json.NewEncoder(w).Encode(UserHandlerOutputDTO{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
