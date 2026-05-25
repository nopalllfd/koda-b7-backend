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
	"math"
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

func (ts *TransactionService) GetAllUserTransaction(ctx context.Context, userID int, query dto.TransactionQuery) (*dto.TransactionPaginationResponse, error) {
	transactions, total, err := ts.transactionRepo.GetAllByUserId(
		ctx,
		ts.db,
		userID,
		dto.TransactionQuery{
			Page:   query.Page,
			Limit:  query.Limit,
			Search: query.Search,
		},
	)

	if err != nil {
		return nil, err
	}

	var response []dto.TransactionResponse

	for _, trx := range transactions {

		item := dto.TransactionResponse{
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

		response = append(response, item)
	}

	// default limit protection
	log.Println(query)
	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	result := &dto.TransactionPaginationResponse{
		Data: response,
		Meta: dto.PaginationMeta{
			Page:       query.Page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return result, nil
}

func (ts *TransactionService) CreateTopup(ctx context.Context, userID int, input dto.TopupRequest) (dto.TopupResponse, error) {

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return dto.TopupResponse{}, errs.ErrTransactionFailed
	}
	defer tx.Rollback(ctx)

	walletID, err := ts.transactionRepo.GetWalletIdByUserID(ctx, tx, userID)
	if err != nil {
		return dto.TopupResponse{}, err
	}

	adminFee := 2500.00
	taxAmount := adminFee * 0.12
	totalAmount := float64(input.Amount) + taxAmount + adminFee

	timestamp := time.Now().Format("20060102150405")
	refCode := fmt.Sprintf("TOP-%s-%d", timestamp, walletID)
	status := "success"
	paymentRef := fmt.Sprintf("80777%s%d", time.Now().Format("150405"), walletID)

	transactionID, err := ts.transactionRepo.CreateTransaction(ctx, tx, "topup", refCode, status)
	if err != nil {
		return dto.TopupResponse{}, err
	}

	if err := ts.transactionRepo.CreateTopup(ctx, tx, transactionID, walletID, input.MethodID, float64(input.Amount), adminFee, totalAmount, paymentRef); err != nil {
		return dto.TopupResponse{}, err
	}

	currentBalance, err := ts.transactionRepo.GetWalletBalance(ctx, tx, walletID)
	if err != nil {
		return dto.TopupResponse{}, errs.ErrWalletNotFound
	}

	newBalance := currentBalance + float64(input.Amount)
	if err := ts.transactionRepo.UpdateWalletBalance(ctx, tx, walletID, newBalance); err != nil {
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
		Total:            totalAmount,
		Status:           status,
		CreatedAt:        time.Now(),
	}, nil
}

func (ts *TransactionService) CreateTransfer(ctx context.Context, input dto.TransferRequest, userID int) (dto.TransferResponse, error) {

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	existingPin, err := ts.transactionRepo.GetPinByUserId(ctx, ts.db, userID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) || errors.Is(err, errs.ErrPINNotSet) {
			return dto.TransferResponse{}, err
		}
		return dto.TransferResponse{}, errs.ErrInternalServer
	}

	log.Println("ini pin awal", input.Pin)
	if err := hc.Compare(input.Pin, existingPin); err != nil {
		return dto.TransferResponse{}, errs.ErrInvalidPin
	}

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return dto.TransferResponse{}, errs.ErrTransactionFailed
	}
	defer tx.Rollback(ctx)

	SenderWalletID, err := ts.transactionRepo.GetWalletIdByUserID(ctx, tx, userID)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	if SenderWalletID == input.ReceiverWalletID {
		return dto.TransferResponse{}, errs.ErrSameWalletTransfer
	}

	senderBalance, err := ts.transactionRepo.GetWalletBalance(ctx, tx, SenderWalletID)
	if err != nil {
		return dto.TransferResponse{}, errs.ErrWalletNotFound
	}

	if senderBalance < float64(input.Amount) {
		return dto.TransferResponse{}, errs.ErrInsufficientBalance
	}

	receiverBalance, err := ts.transactionRepo.GetWalletBalance(ctx, tx, input.ReceiverWalletID)
	if err != nil {
		return dto.TransferResponse{}, errs.ErrWalletNotFound
	}

	timestamp := time.Now().Format("20060102150405")
	refCode := fmt.Sprintf("TRF-%s-%d", timestamp, SenderWalletID)
	status := "success"

	transactionID, err := ts.transactionRepo.CreateTransaction(ctx, tx, "transfer", refCode, status)
	if err != nil {
		return dto.TransferResponse{}, errs.ErrTransactionFailed
	}

	if err := ts.transactionRepo.CreateTransfer(ctx, tx, transactionID, SenderWalletID, input.ReceiverWalletID, float64(input.Amount), input.Description); err != nil {
		return dto.TransferResponse{}, errs.ErrTransferFailed
	}

	newSenderBalance := senderBalance - float64(input.Amount)
	if err := ts.transactionRepo.UpdateWalletBalance(ctx, tx, SenderWalletID, newSenderBalance); err != nil {
		return dto.TransferResponse{}, errs.ErrUpdateBalanceFailed
	}

	newReceiverBalance := receiverBalance + float64(input.Amount)
	if err := ts.transactionRepo.UpdateWalletBalance(ctx, tx, input.ReceiverWalletID, newReceiverBalance); err != nil {
		return dto.TransferResponse{}, errs.ErrUpdateBalanceFailed
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.TransferResponse{}, errs.ErrTransactionFailed
	}

	return dto.TransferResponse{
		TransactionID:    transactionID,
		ReferenceCode:    refCode,
		SenderWalletID:   SenderWalletID,
		ReceiverWalletID: input.ReceiverWalletID,
		Amount:           float64(input.Amount),
		Description:      input.Description,
		Status:           status,
		CreatedAt:        time.Now(),
	}, nil
}

func (ts *TransactionService) GetPaymentMethods(ctx context.Context) ([]dto.PaymentMethods, error) {
	methods, err := ts.transactionRepo.GetAllPaymentMethods(ctx, ts.db)
	if err != nil {
		return nil, err
	}

	var response []dto.PaymentMethods

	for _, item := range methods {
		method := dto.PaymentMethods{
			Id:   item.Id,
			Name: item.Name,
			Logo: item.Logo,
		}
		response = append(response, method)
	}
	return response, nil
}

func (ts *TransactionService) GetAllReceivers(
	ctx context.Context,
	query dto.TransactionQuery,
	userID int,
) (*dto.ReceiverPaginationResponse, error) {

	data, total, err := ts.transactionRepo.GetReceiversWithPagination(
		ctx,
		ts.db,
		query,
		userID,
	)

	if err != nil {
		return nil, errs.ErrInternalServer
	}

	if len(data) == 0 {
		return nil, errs.ErrNoReceiverFound
	}

	var response []dto.Receivers

	for _, item := range data {
		response = append(response, dto.Receivers{
			Id:       item.Id,
			Photo:    item.Photo,
			FullName: item.FullName,
			Phone:    item.Phone,
		})
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}

	if query.Page <= 0 {
		query.Page = 1
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	result := dto.ReceiverPaginationResponse{
		Data: response,
		Meta: dto.PaginationMeta{
			Page:       query.Page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return &result, nil
}

func (ts *TransactionService) GetChartData(
	ctx context.Context,
	userID int,
	query dto.ChartQuery,
) ([]dto.IncomeExpenseChart, error) {

	var interval string

	if query.Period != "7d" && query.Period != "1m" {
		return nil, errors.New("invalid date range")
	}
	if query.Period == "1m" {
		interval = "30 days"
	} else {
		interval = "7 days"
	}

	txType := query.Type
	if txType == "" {
		txType = "all"
	}

	if txType == "income" {
		txType = "in"
	} else if txType == "expense" {
		txType = "out"
	}

	if txType != "all" && txType != "in" && txType != "out" {
		return nil, errors.New("invalid flow type")
	}

	data, err := ts.transactionRepo.GetChartData(
		ctx,
		ts.db,
		userID,
		interval,
		txType,
	)

	if err != nil {
		return nil, err
	}

	var response []dto.IncomeExpenseChart

	for _, item := range data {

		response = append(response, dto.IncomeExpenseChart{
			Date:   item.Date,
			Amount: item.Amount,
			Type:   item.Type,
		})
	}

	return response, nil
}
