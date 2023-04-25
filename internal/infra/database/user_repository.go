package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/sesaquecruz/go-auth-api/internal/entity"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Save(ctx context.Context, user entity.User) error {
	stmt, err := r.DB.PrepareContext(ctx, "INSERT INTO users (id, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.ID, user.Email, user.Password)
	return err
}

func (r *UserRepository) FindById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, email, password FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user entity.User
	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, email, password FROM users WHERE email = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user entity.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user entity.User) error {
	stmt, err := r.DB.Prepare("UPDATE users SET email = ?, password = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.Email, user.Password, user.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	stmt, err := r.DB.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	return err
}
