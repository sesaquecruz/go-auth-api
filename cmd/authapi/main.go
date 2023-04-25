package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/sesaquecruz/go-auth-api/config"
	"github.com/sesaquecruz/go-auth-api/internal/entity"
	"github.com/sesaquecruz/go-auth-api/internal/infra/database"
	"github.com/sesaquecruz/go-auth-api/internal/infra/web/handler"
	"github.com/sesaquecruz/go-auth-api/internal/usecase"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

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

	userFactory := entity.NewUserFactory()
	userRepository := database.NewUserRepository(db)
	createUserUseCase := usecase.NewCreateUserUseCase(userFactory, userRepository)
	userHandler := handler.NewUserHandler(createUserUseCase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
	})

	log.Println("server is running on port 8080...")
	http.ListenAndServe(":8080", r)
}
