package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// import (
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepo(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (wr *WalletRepository) Create(ctx context.Context, userID, balance int) error {
	sql := "INSERT into wallets (user_id, balance) VALUES ($1,$2)"
	if _, err := wr.db.Exec(ctx, sql, userID, balance); err != nil {
		return err
	}
	return nil
}

// func (wr *WalletRepository) GetDashboard(ctx context.Context) (model.WalletSummary, error) {
// 	sql :=
// }
