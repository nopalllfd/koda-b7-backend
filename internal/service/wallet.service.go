package service

import (
	"backend-golang/internal/dto"
	"backend-golang/internal/repository"
	"context"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
}

func NewWalletService(walletRepo *repository.WalletRepository) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
	}
}

func (ws *WalletService) GetDashboard(ctx context.Context, userID int) (dto.DashboardUser, error) {
	data, err := ws.walletRepo.GetDashboard(ctx, userID)
	if err != nil {
		return dto.DashboardUser{}, err
	}
	dataDashboard := dto.DashboardUser{
		Balance: data.Balance,
		Income:  data.Income,
		Expense: data.Expense,
	}
	return dataDashboard, nil
}
