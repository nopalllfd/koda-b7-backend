package dto

import "time"

type TransactionResponse struct {
	TransactionID     int       `json:"transaction_id" db:"transaction_id"`
	ReferenceCode     string    `json:"reference_code" db:"reference_code"`
	TransactionType   string    `json:"transaction_type" db:"transaction_type"`   // "topup" atau "transfer"
	TransactionLabel  string    `json:"transaction_label" db:"transaction_label"` // "Top Up", "Transfer Masuk"
	FlowType          string    `json:"flow_type" db:"flow_type"`                 // "in" atau "out"
	Amount            float64   `json:"amount" db:"amount"`
	CounterpartyName  *string   `json:"counterparty_name" db:"counterparty_name"`   // Nama lawan / nama payment method
	CounterpartyPhone *string   `json:"counterparty_phone" db:"counterparty_phone"` // Pakai pointer agar aman dari NULL
	Status            string    `json:"status" db:"status"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

type TransactionPaginationResponse struct {
	Data []TransactionResponse `json:"data"`
	Meta PaginationMeta        `json:"meta"`
}

type TopupRequest struct {
	WalletID int `json:"wallet_id"`
	MethodID int `json:"method_id"`
	Amount   int `json:"amount"`
}

type TopupResponse struct {
	TransactionID    int       `json:"transaction_id"`
	ReferenceCode    string    `json:"reference_code"`
	PaymentReference string    `json:"payment_reference"`
	Amount           float64   `json:"amount"`
	AdminFee         float64   `json:"admin_fee"`
	TaxAmount        float64   `json:"tax_amount"`
	Total            float64   `json:"total"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type TransferRequest struct {
	SenderWalletID   int     `json:"sender_wallet_id"`
	ReceiverWalletID int     `json:"receiver_wallet_id"`
	Amount           float64 `json:"amount"`
	Description      string  `json:"description"`
}

type TransferResponse struct {
	TransactionID    int       `json:"transaction_id"`
	ReferenceCode    string    `json:"reference_code"`
	SenderWalletID   int       `json:"sender_wallet_id"`
	ReceiverWalletID int       `json:"receiver_wallet_id"`
	Amount           float64   `json:"amount"`
	Description      string    `json:"description"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type TransactionQuery struct {
	Page   int
	Limit  int
	Search string
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
