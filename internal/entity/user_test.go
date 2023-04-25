package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func Test_User_NewUserFactory(t *testing.T) {
	userFactory := NewUserFactory()
	assert.NotNil(t, userFactory)
}

func Test_User_NewUser(t *testing.T) {
	userFactory := UserFactory{}
	email := "user@mail.com"
	password := "12345"

	user, err := userFactory.NewUser(email, password)
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.NotEqual(t, user.ID, uuid.Nil)
	assert.Equal(t, user.Email, email)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)))

	user, err = userFactory.NewUser("user@mailcom", password)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrUserInvalidEmail)

	user, err = userFactory.NewUser("usermail.com", password)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrUserInvalidEmail)

	user, err = userFactory.NewUser(email, "1234")
	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrUserInvalidPassword)
}

func Test_User_Validate(t *testing.T) {
	user := User{ID: uuid.Nil}
	assert.ErrorIs(t, user.Validate(), ErrUserInvalidID)

	user = User{ID: uuid.New(), Email: "user@mailcom"}
	assert.ErrorIs(t, user.Validate(), ErrUserInvalidEmail)

	user = User{ID: uuid.New(), Email: "user@mail.com", Password: "12345"}
	assert.ErrorIs(t, user.Validate(), ErrUserInvalidPassword)

	hash, _ := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	user = User{ID: uuid.New(), Email: "user@mail.com", Password: string(hash)}
	assert.Nil(t, user.Validate())
}
