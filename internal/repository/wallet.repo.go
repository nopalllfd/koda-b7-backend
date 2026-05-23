package repository

import (
	"backend-golang/internal/model"
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

func (wr *WalletRepository) GetDashboard(ctx context.Context, userID int) (model.WalletSummary, error) {
	sql := `SELECT 
    (
        SELECT COALESCE(SUM(t.amount), 0)
        FROM transfers t
        JOIN wallets w ON w.id = t.receiver_wallet_id
        JOIN transactions trx ON trx.id = t.transaction_id
        WHERE w.user_id = $1 AND trx.status = 'success'
    ) 

	(SELECT COALESCE(SUM(t.amount), 0)
        FROM transfers t
        JOIN wallets w ON w.id = t.sender_wallet_id
        JOIN transactions trx ON trx.id = t.transaction_id
        WHERE w.user_id = $1 AND trx.status = 'success') AS grand_total_expense, wallets.balance FROM wallets WHERE wallets.user_id = $1
	`
	var userDashboard model.WalletSummary
	if err := wr.db.QueryRow(ctx, sql, userID).Scan(&userDashboard.Income, &userDashboard.Expense, &userDashboard.Balance); err != nil {
		return model.WalletSummary{}, err
	}

	return userDashboard, nil

}
