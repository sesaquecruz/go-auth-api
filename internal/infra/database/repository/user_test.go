package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/sesaquecruz/go-auth-api/config"
	"github.com/sesaquecruz/go-auth-api/internal/entity"
	"github.com/stretchr/testify/suite"

	_ "github.com/go-sql-driver/mysql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db             *sql.DB
	migrate        *migrate.Migrate
	userRepository *UserRepository
	ctx            context.Context
	user1          *entity.User
	user2          *entity.User
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	cfg, err := config.LoadConfig()
	s.Require().Nil(err)

	db_url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open(cfg.DBDriver, db_url)
	s.Require().Nil(err)

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	s.Require().Nil(err)

	dir, err := os.Getwd()
	s.Require().Nil(err)

	for i := 0; i < 4; i++ {
		dir = filepath.Dir(dir)
	}

	migrate, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/migrations", dir), cfg.DBName, driver)
	s.Require().Nil(err)

	err = migrate.Up()
	s.Require().Nil(err)

	s.db = db
	s.migrate = migrate
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	err := s.migrate.Down()
	s.Require().Nil(err)

	s.migrate.Close()
	s.db.Close()
}

func (s *UserRepositoryTestSuite) SetupTest() {
	s.userRepository = &UserRepository{DB: s.db}
	s.ctx = context.Background()
	s.user1 = &entity.User{ID: uuid.New(), Email: "user1@mail.com", Password: "12345"}
	s.user2 = &entity.User{ID: uuid.New(), Email: "user2@mail.com", Password: "54321"}
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	_, err := s.db.Exec("DELETE FROM users")
	s.Require().Nil(err)
}

func TestSuite_UserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) Test_UserRepository_NewUserRepository() {
	userRepository := NewUserRepository(s.db)
	s.NotNil(userRepository)
	s.Equal(s.userRepository, userRepository)
}

func (s *UserRepositoryTestSuite) Test_UserRepository_Save() {
	user, err := s.userRepository.FindById(s.ctx, s.user1.ID)
	s.ErrorIs(err, sql.ErrNoRows)
	s.Nil(user)

	err = s.userRepository.Save(s.ctx, *s.user1)
	s.Nil(err)

	user, err = s.userRepository.FindById(s.ctx, s.user1.ID)
	s.Nil(err)
	s.Equal(s.user1, user)
}

func (s *UserRepositoryTestSuite) Test_UserRepository_FindById() {
	user, err := s.userRepository.FindById(s.ctx, s.user1.ID)
	s.ErrorIs(err, sql.ErrNoRows)
	s.Nil(user)

	err = s.userRepository.Save(s.ctx, *s.user1)
	s.Nil(err)

	user, err = s.userRepository.FindById(s.ctx, s.user1.ID)
	s.Nil(err)
	s.Equal(s.user1, user)
}

func (s *UserRepositoryTestSuite) Test_UserRepository_FindByEmail() {
	user, err := s.userRepository.FindByEmail(s.ctx, s.user1.Email)
	s.ErrorIs(err, sql.ErrNoRows)
	s.Nil(user)

	err = s.userRepository.Save(s.ctx, *s.user1)
	s.Nil(err)

	user, err = s.userRepository.FindByEmail(s.ctx, s.user1.Email)
	s.Nil(err)
	s.Equal(s.user1, user)
}

func (s *UserRepositoryTestSuite) Test_UserRepository_Update() {
	err := s.userRepository.Save(s.ctx, *s.user1)
	s.Nil(err)

	user, err := s.userRepository.FindById(s.ctx, s.user1.ID)
	s.Nil(err)
	s.Equal(s.user1, user)
	s.NotEqual(s.user2.Email, user.Email)
	s.NotEqual(s.user2.Password, user.Password)

	err = s.userRepository.Update(s.ctx, entity.User{ID: s.user1.ID, Email: s.user2.Email, Password: s.user2.Password})
	s.Nil(err)

	user, err = s.userRepository.FindById(s.ctx, s.user1.ID)
	s.Nil(err)
	s.Equal(s.user1.ID, user.ID)
	s.Equal(s.user2.Email, user.Email)
	s.Equal(s.user2.Password, user.Password)
}

func (s *UserRepositoryTestSuite) Test_UserRepository_Delete() {
	err := s.userRepository.Save(s.ctx, *s.user1)
	s.Nil(err)

	err = s.userRepository.Save(s.ctx, *s.user2)
	s.Nil(err)

	user, err := s.userRepository.FindById(s.ctx, s.user1.ID)
	s.Nil(err)
	s.Equal(s.user1, user)

	user, err = s.userRepository.FindById(s.ctx, s.user2.ID)
	s.Nil(err)
	s.Equal(s.user2, user)

	err = s.userRepository.Delete(s.ctx, s.user1.ID)
	s.Nil(err)

	user, err = s.userRepository.FindById(s.ctx, s.user1.ID)
	s.ErrorIs(err, sql.ErrNoRows)
	s.Nil(user)

	err = s.userRepository.Delete(s.ctx, s.user2.ID)
	s.Nil(err)

	user, err = s.userRepository.FindById(s.ctx, s.user2.ID)
	s.ErrorIs(err, sql.ErrNoRows)
	s.Nil(user)
}
