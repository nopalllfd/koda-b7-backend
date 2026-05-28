package repository

import (
	"backend-golang/internal/dto"
	errs "backend-golang/internal/err"
	"backend-golang/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

type DBTX interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type TransactionRepository struct {
	rc *redis.Client
}

func NewTransactionRepo(rc *redis.Client) *TransactionRepository {
	return &TransactionRepository{
		rc: rc,
	}
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
	query dto.TransactionQuery,
) ([]model.TransactionResponse, int64, error) {

	// =========================
	// DEFAULT PAGINATION
	// =========================

	if query.Page <= 0 {
		query.Page = 1
	}

	if query.Limit <= 0 {
		query.Limit = 10
	}

	if query.Limit > 100 {
		query.Limit = 100
	}

	offset := (query.Page - 1) * query.Limit

	// =========================
	// CACHE KEY
	// =========================

	key := fmt.Sprintf(
		"trx:user:%d:page:%d:limit:%d",
		id,
		query.Page,
		query.Limit,
	)

	// =========================
	// GET CACHE
	// =========================

	cached, err := tr.rc.Get(
		ctx,
		key,
	).Result()

	if err == nil {

		var result struct {
			Transactions []model.TransactionResponse `json:"transactions"`
			Total        int64                       `json:"total"`
		}

		if err := json.Unmarshal(
			[]byte(cached),
			&result,
		); err == nil {

			return result.Transactions, result.Total, nil
		}
	}

	// =========================
	// MAIN QUERY
	// =========================

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

	ORDER BY trx.created_at DESC

	LIMIT $2
	OFFSET $3
	`

	rows, err := dbtx.Query(
		ctx,
		sql,
		id,
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

		transactions = append(
			transactions,
			item,
		)
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
		SELECT transaction_id
		FROM topups
		WHERE wallet_id = (SELECT id FROM MyWallet)
	),

	MyTransfers AS (
		SELECT transaction_id
		FROM transfers
		WHERE 
			sender_wallet_id = (SELECT id FROM MyWallet)
			OR receiver_wallet_id = (SELECT id FROM MyWallet)
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
	`

	var total int64

	err = dbtx.QueryRow(
		ctx,
		countSQL,
		id,
	).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	// =========================
	// SET CACHE
	// =========================

	cacheData := struct {
		Transactions []model.TransactionResponse `json:"transactions"`
		Total        int64                       `json:"total"`
	}{
		Transactions: transactions,
		Total:        total,
	}

	jsonData, err := json.Marshal(cacheData)

	if err == nil {

		tr.rc.Set(
			ctx,
			key,
			jsonData,
			30*time.Second,
		)
	}

	return transactions, total, nil
}

func (tr *TransactionRepository) CreateTransaction(ctx context.Context, dbtx DBTX, txType, refCode, status string) (int, error) {
	sql := `
		INSERT INTO transactions (type, reference_code, status, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id;
	`
	var trxID int
	if err := dbtx.QueryRow(ctx, sql, txType, refCode, status).Scan(&trxID); err != nil {
		return 0, err
	}
	return trxID, nil
}

func (tr *TransactionRepository) GetWalletIdByUserID(ctx context.Context, dbtx DBTX, userID int) (int, error) {
	sql := `SELECT id FROM wallets WHERE user_id = $1`
	var id int
	if err := dbtx.QueryRow(ctx, sql, userID).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (tr *TransactionRepository) CreateTopup(ctx context.Context, dbtx DBTX, transactionID, walletID, methodID int, amount, adminFee, total float64, paymentRef string) error {
	query := `
		INSERT INTO topups (transaction_id, wallet_id, method_id, amount, admin_fee, total, payment_reference)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	if _, err := dbtx.Exec(ctx, query, transactionID, walletID, methodID, amount, adminFee, total, paymentRef); err != nil {
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

func (tr *TransactionRepository) GetAllPaymentMethods(ctx context.Context, dbtx DBTX) ([]model.PaymentMethods, error) {
	sql := `SELECT id, name, logo, created_at, updated_at FROM payment_methods`
	rows, err := dbtx.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var methods []model.PaymentMethods

	for rows.Next() {
		var item model.PaymentMethods

		if err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Logo,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		methods = append(methods, item)

	}

	return methods, nil
}

func (ts *TransactionRepository) GetReceiversWithPagination(ctx context.Context, dbtx DBTX, query dto.TransactionQuery, userID int) ([]model.Receivers, int64, error) {
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
		SELECT u.id, p.photo, p.full_name, p.phone
	FROM users u
	JOIN profiles p ON p.user_id = u.id
	WHERE p.phone IS NOT NULL
	AND p.phone <> ''
	AND p.photo IS NOT NULL
	AND p.photo <> ''
	AND (
	p.phone ILIKE '%' || $1 || '%'
	OR
	p.full_name ILIKE '%' || $1 || '%'
	)
	AND u.id <> $4
	ORDER BY u.id DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := dbtx.Query(ctx, sql, query.Search, query.Limit, offset, userID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []model.Receivers
	for rows.Next() {
		var item model.Receivers

		if err := rows.Scan(
			&item.Id,
			&item.Photo,
			&item.FullName,
			&item.Phone,
		); err != nil {
			return nil, 0, err
		}
		if err := rows.Err(); err != nil {
			return nil, 0, err
		}
		result = append(result, item)
	}

	countQuery := `
	SELECT COUNT(*)
FROM users u
JOIN profiles p ON p.user_id = u.id
WHERE p.phone IS NOT NULL
AND p.phone <> ''
AND p.photo IS NOT NULL
AND p.photo <> ''
AND (
	p.phone ILIKE '%' || $1 || '%'
	OR
	p.full_name ILIKE '%' || $1 || '%'
)
		AND u.id <> $2
`

	var total int64

	if err := dbtx.QueryRow(ctx, countQuery, query.Search, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	return result, total, nil

}

func (tr *TransactionRepository) GetChartData(
	ctx context.Context,
	dbtx DBTX,
	userID int,
	interval string,
	txType string,
) ([]model.IncomeExpenseChart, error) {

	sql := `
WITH MyWallet AS (
	SELECT id
	FROM wallets
	WHERE user_id = $1
	LIMIT 1
),

MyTransfers AS (
	SELECT 
		tf.transaction_id,
		tf.amount,
		CASE 
			WHEN tf.receiver_wallet_id = (SELECT id FROM MyWallet) THEN 'in'
			ELSE 'out'
		END AS flow,
		trx.created_at::date AS trx_date
	FROM transfers tf
	JOIN transactions trx ON trx.id = tf.transaction_id
	WHERE tf.sender_wallet_id = (SELECT id FROM MyWallet)
	   OR tf.receiver_wallet_id = (SELECT id FROM MyWallet)
),

AllTx AS (
	SELECT * FROM MyTransfers
)

SELECT 
	trx_date,
	SUM(amount) AS amount,
	flow
FROM AllTx
WHERE trx_date >= (CURRENT_DATE - $2::interval)
  AND ($3 = 'all' OR flow = $3)
GROUP BY trx_date, flow
ORDER BY trx_date ASC;
`

	rows, err := dbtx.Query(ctx, sql, userID, interval, txType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.IncomeExpenseChart

	for rows.Next() {
		var item model.IncomeExpenseChart

		if err := rows.Scan(&item.Date, &item.Amount, &item.Type); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

func (tr *TransactionRepository) DeleteTransactionCache(
	ctx context.Context,
	userID int,
) error {

	keys, err := tr.rc.Keys(
		ctx,
		fmt.Sprintf(
			"trx:user:%d:*",
			userID,
		),
	).Result()

	if err != nil {
		return err
	}

	if len(keys) > 0 {

		return tr.rc.Del(
			ctx,
			keys...,
		).Err()
	}

	return nil
}
