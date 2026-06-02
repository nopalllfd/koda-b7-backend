package model

import "time"

type TransactionResponse struct {
	TransactionID     int       `json:"transaction_id"`
	ReferenceCode     string    `json:"reference_code"`
	TransactionType   string    `json:"transaction_type"`
	TransactionLabel  string    `json:"transaction_label"`
	FlowType          string    `json:"flow_type"`
	Amount            float64   `json:"amount"`
	CounterpartyName  *string   `json:"counterparty_name"`
	CounterpartyPhone *string   `json:"counterparty_phone"`
	Photo             *string   `json:"photo"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
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
