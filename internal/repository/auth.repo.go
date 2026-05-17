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
	sql := `SELECT id, email, password
	FROM users
	WHERE email = $1`
	// ngeeksekusi query
	var user model.User
	err := ar.db.QueryRow(ctx, sql, email).Scan(&user.Id, &user.Email, &user.Password)
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

func (ar *AuthRepository) CreatePin(ctx context.Context, pin string, id int) error {
	sql := `UPDATE users
	SET pin = $1
	WHERE id = $2
	`
	if _, err := ar.db.Exec(ctx, sql, pin, id); err != nil {
		return err
	}
	return nil
}
