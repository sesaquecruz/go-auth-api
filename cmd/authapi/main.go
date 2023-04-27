package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sesaquecruz/go-auth-api/config"
	"github.com/sesaquecruz/go-auth-api/internal/entity"
	"github.com/sesaquecruz/go-auth-api/internal/infra/database/repository"
	"github.com/sesaquecruz/go-auth-api/internal/infra/web/handler"
	mw "github.com/sesaquecruz/go-auth-api/internal/infra/web/middleware"
	"github.com/sesaquecruz/go-auth-api/internal/usecase"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(
		cfg.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName),
	)
	if err != nil {
		panic(err)
	}

	jwtAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	jwtExpiration := time.Duration(cfg.JWTExpSeconds) * time.Second

	userFactory := entity.NewUserFactory()
	userRepository := repository.NewUserRepository(db)

	createUserUseCase := usecase.NewCreateUserUseCase(userFactory, userRepository)
	authUserUseCase := usecase.NewAuthUserUseCase(userFactory, userRepository)
	updateUserUseCase := usecase.NewUpdateUserUseCase(userFactory, userRepository)
	deleteUserUseCase := usecase.NewDeleteUserUseCase(userRepository)

	userHandler := handler.NewUserHandler(
		jwtAuth,
		jwtExpiration,
		createUserUseCase,
		authUserUseCase,
		updateUserUseCase,
		deleteUserUseCase,
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/login", func(r chi.Router) {
		r.Post("/new", userHandler.CreateUser)
		r.Post("/", userHandler.AuthUser)
	})

	r.Route("/user", func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(mw.EchoAuthToken)

		r.Put("/", userHandler.UpdateUser)
		r.Delete("/", userHandler.DeleteUser)
	})

	log.Println("server is running on port 8080...")
	http.ListenAndServe(":8080", r)
}
