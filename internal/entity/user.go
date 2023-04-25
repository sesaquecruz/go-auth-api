package entity

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserInvalidID       = errors.New("invalid id")
	ErrUserInvalidEmail    = errors.New("invalid email")
	ErrUserInvalidPassword = errors.New("invalid password")

	emailPattern    = regexp.MustCompile(`^[a-zA-Z0-9._]+@[a-zA-Z0-9.-]+?\.[a-zA-Z]{2,}$`)
	passwordPattern = regexp.MustCompile(`\$2[ayb]\$.{56}$`)
)

const passwordMinLen = 5

type UserFactory struct{}

func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

func (f *UserFactory) NewUser(email string, password string) (*User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	if !emailPattern.MatchString(email) {
		return nil, ErrUserInvalidEmail
	}

	if len(password) < passwordMinLen {
		return nil, ErrUserInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       id,
		Email:    email,
		Password: string(hash),
	}

	return user, nil
}

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
}

func (u *User) Validate() error {
	if u.ID == uuid.Nil {
		return ErrUserInvalidID
	}
	if !emailPattern.MatchString(u.Email) {
		return ErrUserInvalidEmail
	}
	if !passwordPattern.MatchString(u.Password) {
		return ErrUserInvalidPassword
	}
	return nil
}
