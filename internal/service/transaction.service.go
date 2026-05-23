package service

import (
	"backend-golang/internal/dto"
	errs "backend-golang/internal/err"
	"backend-golang/internal/repository"
	"backend-golang/pkg"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	db              *pgxpool.Pool
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, db *pgxpool.Pool) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		db:              db,
	}
}

func (ts *TransactionService) CheckPin(ctx context.Context, userID int, pin string) error {
	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()
	existingPin, err := ts.transactionRepo.GetPinByUserId(ctx, ts.db, userID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) || errors.Is(err, errs.ErrPINNotSet) {
			return err
		}
		return errs.ErrInternalServer
	}

	log.Println("ini pin awal", pin)
	if err := hc.Compare(pin, existingPin); err != nil {
		return errs.ErrInvalidPin
	}
	return nil
}

func (ts *TransactionService) GetAllUserTransaction(ctx context.Context, userID int) ([]dto.TransactionResponse, error) {
	transactions, err := ts.transactionRepo.GetAllByUserId(ctx, ts.db, userID)
	if err != nil {
		return nil, err
	}

	var response []dto.TransactionResponse

	for _, trx := range transactions {
		items := dto.TransactionResponse{
			TransactionID:     trx.TransactionID,
			ReferenceCode:     trx.ReferenceCode,
			TransactionType:   trx.TransactionType,
			TransactionLabel:  trx.TransactionLabel,
			FlowType:          trx.FlowType,
			Amount:            trx.Amount,
			CounterpartyName:  trx.CounterpartyName,
			CounterpartyPhone: trx.CounterpartyPhone,
			Status:            trx.Status,
			CreatedAt:         trx.CreatedAt,
		}
		response = append(response, items)
	}
	return response, nil
}

func (ts *TransactionService) CreateTopup(ctx context.Context, input dto.TopupRequest) (dto.TopupResponse, error) {

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return dto.TopupResponse{}, errs.ErrTransactionFailed
	}
	defer tx.Rollback(ctx)

	adminFee := 2500.00
	taxAmount := adminFee * 0.12
	totalAmount := float64(input.Amount) + taxAmount + adminFee

	timestamp := time.Now().Format("20060102150405")
	refCode := fmt.Sprintf("TOP-%s-%d", timestamp, input.WalletID)
	status := "success"
	paymentRef := fmt.Sprintf("80777%s%d", time.Now().Format("150405"), input.WalletID)

	transactionID, err := ts.transactionRepo.CreateTransaction(ctx, tx, "topup", refCode, status)
	if err != nil {
		return dto.TopupResponse{}, errs.ErrTransactionFailed
	}

	if err := ts.transactionRepo.CreateTopup(ctx, tx, transactionID, input.WalletID, input.MethodID, float64(input.Amount), adminFee, taxAmount, totalAmount, paymentRef); err != nil {
		return dto.TopupResponse{}, errs.ErrTopupFailed
	}

	currentBalance, err := ts.transactionRepo.GetWalletBalance(ctx, tx, input.WalletID)
	if err != nil {
		return dto.TopupResponse{}, errs.ErrWalletNotFound
	}

	newBalance := currentBalance + float64(input.Amount)
	if err := ts.transactionRepo.UpdateWalletBalance(ctx, tx, input.WalletID, newBalance); err != nil {
		return dto.TopupResponse{}, errs.ErrUpdateBalanceFailed
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.TopupResponse{}, errs.ErrTransactionFailed
	}

	return dto.TopupResponse{
		TransactionID:    transactionID,
		ReferenceCode:    refCode,
		PaymentReference: paymentRef,
		Amount:           float64(input.Amount),
		AdminFee:         adminFee,
		TaxAmount:        taxAmount,
		Total:            totalAmount,
		Status:           status,
		CreatedAt:        time.Now(),
	}, nil
}
func (tx *TransactionService) CreateTransfer(ctx context.Context, input dto.TransferRequest) (dto.TransferResponse, error) {

}
