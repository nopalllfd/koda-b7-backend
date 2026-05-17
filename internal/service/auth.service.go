package service

import (
	"backend-golang/internal/dto"
	errs "backend-golang/internal/err"
	"backend-golang/internal/repository"
	"backend-golang/pkg"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	authRepo   *repository.AuthRepository
	userRepo   *repository.UserRepository
	walletRepo *repository.WalletRepository
}

func NewAuthService(authRepo *repository.AuthRepository, userRepo *repository.UserRepository, walletRepo *repository.WalletRepository) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

func (as *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	log.Println("MASUK LOGIN SERVICE")
	existingUser, err := as.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrEmailNotFound
		}
		return nil, err
	}
	log.Println("HABIS GET USER BY MAIL")

	var hc pkg.HashConfig

	if err := hc.Compare(req.Password, existingUser.Password); err != nil {
		return nil, errs.ErrInvalidCredential
	}
	log.Println("HABIS COMPARE PW")

	log.Println("SEBELUM GEN JWT")

	claims := pkg.NewClaims(existingUser.Id, req.Email)
	token, err := claims.GenJWT()
	if err != nil {
		return nil, err
	}

	result := &dto.LoginResponse{
		ID:    existingUser.Id,
		Email: existingUser.Email,
		Token: token,
	}

	log.Println("SETELAH GEN JWT:", token)
	return result, nil
}

func (as *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	isEmailExists, err := as.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return errs.ErrInternalServer
	}
	if isEmailExists {
		return errs.ErrExistingEmail
	}
	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()
	hashedPassword := hc.Hash(req.Password)
	log.Println(hashedPassword)
	userID, err := as.authRepo.Create(ctx, req.Email, hashedPassword)
	if err != nil {
		return errs.ErrInternalServer
	}

	//create profile dan create wallet

	log.Println("USER ID:", userID)

	if err := as.userRepo.Create(ctx, userID, nil, nil, nil); err != nil {
		log.Println("PROFILE ERROR:", err)
		return err
	}

	log.Println("PROFILE CREATED")
	if err := as.walletRepo.Create(ctx, userID, 0); err != nil {
		log.Println("WALLET ERROR:", err)
		return err
	}

	return nil
}

func (as *AuthService) AddPin(ctx context.Context, req dto.AddPinRequest) error {
	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()
	hashedPin := hc.Hash(req.Pin)
	if err := as.authRepo.CreatePin(ctx, hashedPin, req.UserID); err != nil {
		return err
	}
	return nil
}
