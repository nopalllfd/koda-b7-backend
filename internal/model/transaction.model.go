package model

import "time"

type TransactionResponse struct {
	TransactionID     int       `db:"transaction_id"`
	ReferenceCode     string    `db:"reference_code"`
	TransactionType   string    `db:"transaction_type"`  // "topup" atau "transfer"
	TransactionLabel  string    `db:"transaction_label"` // "Top Up", "Transfer Masuk"
	FlowType          string    `db:"flow_type"`         // "in" atau "out"
	Amount            float64   `db:"amount"`
	CounterpartyName  *string   `db:"counterparty_name"`  // Nama lawan / nama payment method
	CounterpartyPhone *string   `db:"counterparty_phone"` // Pakai pointer agar aman dari NULL
	Status            string    `db:"status"`
	CreatedAt         time.Time `db:"created_at"`
}

type Topups struct {
	TransactionID    int        `db:"transaction_id"`
	WalletID         int        `db:"wallet_id"`
	MethodID         int        `db:"method_id"`
	Amount           float64    `db:"amount"`
	AdminFee         *float64   `db:"admin_fee"`
	Total            float64    `db:"total"`
	PaymentReference *string    `db:"payment_reference"`
	ExpiredAt        *time.Time `db:"expired_at"`
	PaidAt           *time.Time `db:"paid_at"`
}

type PaymentMethods struct {
	Id        int        `db:"id"`
	Name      string     `db:"name"`
	Logo      string     `db:"logo"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type Receivers struct {
	Id       int    `db:"id"`
	Photo    string `db:"photo"`
	FullName string `db:"full_name"`
	Phone    string `db:"phone"`
}

type IncomeExpenseChart struct {
	Date   time.Time
	Amount float64
	Type   string
}
