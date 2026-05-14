package service

import (
	"backend-golang/internal/dto"
	"backend-golang/internal/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

var ErrInvalidCredential = errors.New(
	"invalid credential",
)

var ErrEmailNotFound = errors.New(
	"email not found",
)
var ErrExistingEmail = errors.New(
	"email has been registered",
)

var ErrInternalServer = errors.New(
	"internal server error",
)

func (as *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := as.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmailNotFound
		}
		return nil, err
	}

	if user == nil {
		return nil, ErrEmailNotFound
	}
	match := req.Password == user.Password
	if !match {
		return nil, ErrInvalidCredential
	}

	token := "INITOKEN1238127381273"

	result := &dto.LoginResponse{
		ID:    user.Id,
		Email: user.Email,
		Token: token,
	}
	return result, nil
}

func (as *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	isEmailExists, err := as.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInternalServer
		}
	}
	if isEmailExists {
		return ErrExistingEmail
	}
	if err := as.authRepo.Create(ctx, req.Email, req.Password); err != nil {
		return ErrInternalServer
	}
	return nil
}
