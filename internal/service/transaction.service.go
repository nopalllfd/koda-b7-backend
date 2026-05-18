package service

import (
	errs "backend-golang/internal/err"
	"backend-golang/internal/repository"
	"backend-golang/pkg"
	"context"
	"log"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

func (ts *TransactionService) CheckPin(ctx context.Context, userID int, pin string) error {
	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()
	existingPin, err := ts.transactionRepo.GetPinByUserId(ctx, userID)
	if err != nil {
		return err
	}

	log.Println("ini pin awal", pin)
	if err := hc.Compare(pin, existingPin); err != nil {
		return errs.ErrInvalidPin
	}
	return nil
}
