package repository

import (
	errs "backend-golang/internal/err"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (tr *TransactionRepository) GetPinByUserId(ctx context.Context, id int) (string, error) {
	sql := `SELECT COALESCE(pin, '') FROM users WHERE id=$1`
	var pin string
	if err := tr.db.QueryRow(ctx, sql, id).Scan(&pin); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errs.ErrUserNotFound
		}
		return "", err
	}
	log.Println("Ini log", pin)
	return pin, nil
}
