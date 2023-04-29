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
	"github.com/sesaquecruz/go-auth-api/internal/usecase"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/sesaquecruz/go-auth-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

const basePath = "/api/v1"
const port = "8080"

// @title          	Auth API
// @version        	1.0.0
// @description    	An Auth API with JWT and RSA

// @contact.name   	API Repository
// @contact.url    	https://github.com/sesaquecruz/go-auth-api

// @license.name	MIT License
// @license.url   	https://github.com/sesaquecruz/go-auth-api/blob/main/LICENSE

// @host           	localhost:8080
// @BasePath       	/api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in             	header
// @name           	Authorization
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
	findUserUseCase := usecase.NewFindUserUseCase(userRepository)

	userHandler := handler.NewUserHandler(
		jwtAuth,
		jwtExpiration,
		createUserUseCase,
		authUserUseCase,
		updateUserUseCase,
		deleteUserUseCase,
		findUserUseCase,
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	authMiddlewares := chi.Chain(
		jwtauth.Verifier(jwtAuth),
		jwtauth.Authenticator,
	)

	r.Route(basePath+"/login", func(r chi.Router) {
		r.Post("/", userHandler.AuthUser)
	})

	r.Route(basePath+"/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.With(authMiddlewares...).Get("/", userHandler.FindUser)
		r.With(authMiddlewares...).Put("/", userHandler.UpdateUser)
		r.With(authMiddlewares...).Delete("/", userHandler.DeleteUser)
	})

	r.Get(
		basePath+"/docs/*",
		httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://localhost:%s%s/docs/doc.json", port, basePath))),
	)

	log.Printf("server is running on port %s...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
