package repository

import (
	errs "backend-golang/internal/err"
	"backend-golang/internal/model"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type TransactionRepository struct {
}

func NewTransactionRepo() *TransactionRepository {
	return &TransactionRepository{}
}

func (tr *TransactionRepository) GetPinByUserId(ctx context.Context, dbtx DBTX, id int) (string, error) {
	sql := `SELECT COALESCE(pin, '') FROM users WHERE id=$1`
	var pin string
	if err := dbtx.QueryRow(ctx, sql, id).Scan(&pin); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errs.ErrUserNotFound
		}
		return "", err
	}
	if pin == "" {
		return "", errs.ErrPINNotSet
	}
	log.Println("Ini log", pin)
	return pin, nil
}

func (tr *TransactionRepository) GetAllByUserId(
	ctx context.Context,
	dbtx DBTX,
	id int,
	query model.TransactionQuery,
) ([]model.TransactionResponse, int64, error) {

	// default pagination
	if query.Page <= 0 {
		query.Page = 1
	}

	if query.Limit <= 0 {
		query.Limit = 10
	}

	// max limit protection
	if query.Limit > 100 {
		query.Limit = 100
	}

	offset := (query.Page - 1) * query.Limit

	sql := `
	WITH MyWallet AS (
		SELECT id
		FROM wallets
		WHERE user_id = $1
		LIMIT 1
	),

	MyTopups AS (
		SELECT 
			tp.transaction_id, 
			tp.amount, 
			'Top Up' AS label, 
			'in' AS flow,
			pm.name AS counterparty_name, 
			NULL AS counterparty_phone 
		FROM topups tp
		LEFT JOIN payment_methods pm 
			ON tp.method_id = pm.id
		WHERE tp.wallet_id = (SELECT id FROM MyWallet)
	),

	MyTransfers AS (
		SELECT 
			tf.transaction_id, 
			tf.amount,

			CASE 
				WHEN tf.receiver_wallet_id = (SELECT id FROM MyWallet)
					THEN 'Transfer Masuk'
				ELSE 'Transfer Keluar'
			END AS label,

			CASE 
				WHEN tf.receiver_wallet_id = (SELECT id FROM MyWallet)
					THEN 'in'
				ELSE 'out'
			END AS flow,

			p.full_name AS counterparty_name,
			p.phone AS counterparty_phone

		FROM transfers tf

		JOIN wallets w_lawan 
			ON w_lawan.id = CASE 
				WHEN tf.receiver_wallet_id = (SELECT id FROM MyWallet)
					THEN tf.sender_wallet_id 
				ELSE tf.receiver_wallet_id 
			END

		JOIN profiles p 
			ON p.user_id = w_lawan.user_id

		WHERE 
			tf.sender_wallet_id = (SELECT id FROM MyWallet)
			OR tf.receiver_wallet_id = (SELECT id FROM MyWallet)
	)

	SELECT 
		trx.id AS transaction_id, 
		trx.reference_code,
		trx.type AS transaction_type,

		COALESCE(tp.label, tf.label) AS transaction_label,
		COALESCE(tp.flow, tf.flow) AS flow_type,
		COALESCE(tp.amount, tf.amount) AS amount,

		COALESCE(
			tp.counterparty_name,
			tf.counterparty_name
		) AS counterparty_name,

		COALESCE(
			tp.counterparty_phone,
			tf.counterparty_phone
		) AS counterparty_phone,

		trx.status, 
		trx.created_at

	FROM transactions trx

	LEFT JOIN MyTopups tp 
		ON trx.id = tp.transaction_id

	LEFT JOIN MyTransfers tf 
		ON trx.id = tf.transaction_id

	WHERE 
		(tp.transaction_id IS NOT NULL 
		OR tf.transaction_id IS NOT NULL)

		AND (
			$2 = ''
			OR trx.reference_code ILIKE '%' || $2 || '%'
			OR COALESCE(
				tp.counterparty_name,
				tf.counterparty_name
			) ILIKE '%' || $2 || '%'
		)

	ORDER BY trx.created_at DESC

	LIMIT $3
	OFFSET $4
	`

	rows, err := dbtx.Query(
		ctx,
		sql,
		id,
		query.Search,
		query.Limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []model.TransactionResponse

	for rows.Next() {
		var item model.TransactionResponse

		err := rows.Scan(
			&item.TransactionID,
			&item.ReferenceCode,
			&item.TransactionType,
			&item.TransactionLabel,
			&item.FlowType,
			&item.Amount,
			&item.CounterpartyName,
			&item.CounterpartyPhone,
			&item.Status,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		transactions = append(transactions, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// =========================
	// COUNT QUERY
	// =========================

	countSQL := `
	WITH MyWallet AS (
		SELECT id
		FROM wallets
		WHERE user_id = $1
		LIMIT 1
	),

	MyTopups AS (
		SELECT 
			tp.transaction_id,
			pm.name AS counterparty_name
		FROM topups tp
		LEFT JOIN payment_methods pm 
			ON tp.method_id = pm.id
		WHERE tp.wallet_id = (SELECT id FROM MyWallet)
	),

	MyTransfers AS (
		SELECT 
			tf.transaction_id,
			p.full_name AS counterparty_name

		FROM transfers tf

		JOIN wallets w_lawan 
			ON w_lawan.id = CASE 
				WHEN tf.receiver_wallet_id = (SELECT id FROM MyWallet)
					THEN tf.sender_wallet_id 
				ELSE tf.receiver_wallet_id 
			END

		JOIN profiles p 
			ON p.user_id = w_lawan.user_id

		WHERE 
			tf.sender_wallet_id = (SELECT id FROM MyWallet)
			OR tf.receiver_wallet_id = (SELECT id FROM MyWallet)
	)

	SELECT COUNT(*)

	FROM transactions trx

	LEFT JOIN MyTopups tp 
		ON trx.id = tp.transaction_id

	LEFT JOIN MyTransfers tf 
		ON trx.id = tf.transaction_id

	WHERE 
		(tp.transaction_id IS NOT NULL 
		OR tf.transaction_id IS NOT NULL)

		AND (
			$2 = ''
			OR trx.reference_code ILIKE '%' || $2 || '%'
			OR COALESCE(
				tp.counterparty_name,
				tf.counterparty_name
			) ILIKE '%' || $2 || '%'
		)
	`

	var total int64

	err = dbtx.QueryRow(
		ctx,
		countSQL,
		id,
		query.Search,
	).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (tr *TransactionRepository) CreateTransaction(ctx context.Context, dbtx DBTX, txType, refCode, status string) (int, error) {
	sql := `
		INSERT INTO transactions (type, reference_code, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id;
	`
	var trxID int
	if err := dbtx.QueryRow(ctx, sql, txType, refCode, status).Scan(&trxID); err != nil {
		return 0, err
	}
	return trxID, nil
}

func (tr *TransactionRepository) CreateTopup(ctx context.Context, dbtx DBTX, transactionID, walletID, methodID int, amount, adminFee, taxAmount, total float64, paymentRef string) error {
	query := `
		INSERT INTO topups (transaction_id, wallet_id, method_id, amount, admin_fee, tax_amount, total, payment_reference)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	if _, err := dbtx.Exec(ctx, query, transactionID, walletID, methodID, amount, adminFee, taxAmount, total, paymentRef); err != nil {
		return err
	}
	return nil
}

func (tr *TransactionRepository) GetWalletBalance(ctx context.Context, dbtx DBTX, walletID int) (float64, error) {
	query := `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE;`
	var balance float64
	if err := dbtx.QueryRow(ctx, query, walletID).Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}

func (tr *TransactionRepository) UpdateWalletBalance(ctx context.Context, dbtx DBTX, walletID int, newBalance float64) error {
	query := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2;`
	if _, err := dbtx.Exec(ctx, query, newBalance, walletID); err != nil {
		return err
	}
	return nil
}

func (tr *TransactionRepository) CreateTransfer(ctx context.Context, dbtx DBTX, transactionID, senderWalletID, receiverWalletID int, amount float64, description string) error {
	query := `
		INSERT INTO transfers (transaction_id, sender_wallet_id, receiver_wallet_id, amount, description)
		VALUES ($1, $2, $3, $4, $5);
	`
	if _, err := dbtx.Exec(ctx, query, transactionID, senderWalletID, receiverWalletID, amount, description); err != nil {
		return err
	}
	return nil
}
