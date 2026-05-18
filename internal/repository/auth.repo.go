package repository

import (
	"backend-golang/internal/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (ar *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	// definisiin
	sql := `SELECT id, email, password, COALESCE(pin, '')
	FROM users
	WHERE email = $1`
	// ngeeksekusi query
	var user model.User
	err := ar.db.QueryRow(ctx, sql, email).Scan(&user.Id, &user.Email, &user.Password, &user.Pin)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf(
			"find user by email: %w",
			err,
		)
	}

	return &user, nil
}

func (ar *AuthRepository) FindByEmail(ctx context.Context, email string) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	if err := ar.db.QueryRow(ctx, sql, email).Scan(&exists); err != nil {
		return false, fmt.Errorf(
			"find user by email: %w",
			err,
		)
	}
	return exists, nil
}

func (ar *AuthRepository) Create(ctx context.Context, email string, password string) (int, error) {
	sql := `INSERT INTO users (email,password) VALUES ($1,$2) RETURNING id`

	var userId int
	err := ar.db.QueryRow(ctx, sql, email, password).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf(
			"create user: %w",
			err,
		)
	}

	return userId, nil
}

func (ar *AuthRepository) SetPin(ctx context.Context, pin string, id int) error {
	sql := `UPDATE users
	SET pin = $1
	WHERE id = $2
	`
	if _, err := ar.db.Exec(ctx, sql, pin, id); err != nil {
		return err
	}
	return nil
}
func (ar *AuthRepository) GetUserPassword(ctx context.Context, id int) (string, error) {
	sql := `SELECT password FROM users WHERE id = $1`
	var userPassword string
	if err := ar.db.QueryRow(ctx, sql, id).Scan(&userPassword); err != nil {
		return "", err
	}

	return userPassword, nil
}

func (ar *AuthRepository) SetPassword(ctx context.Context, newPass string, id int) error {
	sql := `UPDATE users
	SET password = $1
	WHERE id = $2
	`
	if _, err := ar.db.Exec(ctx, sql, newPass, id); err != nil {
		return err
	}

	return nil
}

// func (ar *AuthRepository) GetPin(ctx )
