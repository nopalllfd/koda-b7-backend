package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/nopalllfd/koda-b7-backend/internal/dto"
	errs "github.com/nopalllfd/koda-b7-backend/internal/err"
	"github.com/nopalllfd/koda-b7-backend/internal/repository"
	"github.com/nopalllfd/koda-b7-backend/pkg"
	"github.com/nopalllfd/koda-b7-backend/pkg/utils"

	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	authRepo   *repository.AuthRepository
	userRepo   *repository.UserRepository
	walletRepo *repository.WalletRepository
	db         *repository.DBTX
}

func NewAuthService(authRepo *repository.AuthRepository, userRepo *repository.UserRepository, walletRepo *repository.WalletRepository) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

func (as *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	log.Println("[Login] START")

	existingUser, err := as.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("[Login] GetUserByEmail error: %v\n", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrEmailNotFound
		}

		return nil, errs.ErrInternalServer
	}

	log.Println("[Login] User found")

	var hc pkg.HashConfig

	if err := hc.Compare(req.Password, existingUser.Password); err != nil {
		log.Printf("[Login] Password compare error: %v\n", err)
		return nil, errs.ErrInvalidCredential
	}

	log.Println("[Login] Password valid")

	claims := pkg.NewClaims(existingUser.Id, req.Email)

	token, err := claims.GenJWT()
	if err != nil {
		log.Printf("[Login] Generate JWT error: %v\n", err)
		return nil, errs.ErrInternalServer
	}

	isPinExists := len(existingUser.Pin) > 0

	displayName := existingUser.FullName
	if displayName == "" {
		displayName = existingUser.Email
	}

	result := &dto.LoginResponse{
		DisplayName: displayName,
		Photo:       existingUser.Photo,
		IsPinExists: isPinExists,
		Token:       token,
	}

	log.Println("[Login] SUCCESS")

	return result, nil
}

func (as *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	log.Println("[Register] START")

	isEmailExists, err := as.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("[Register] FindByEmail error: %v\n", err)
		return errs.ErrInternalServer
	}

	if isEmailExists {
		log.Printf("[Register] Email already exists: %s\n", req.Email)
		return errs.ErrExistingEmail
	}

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	hashedPassword := hc.Hash(req.Password)

	userID, err := as.authRepo.Create(ctx, req.Email, hashedPassword)
	if err != nil {
		log.Printf("[Register] Create user error: %v\n", err)
		return errs.ErrInternalServer
	}

	log.Printf("[Register] User created. ID=%d\n", userID)

	if err := as.userRepo.Create(ctx, userID, nil, nil, nil); err != nil {
		log.Printf("[Register] Create profile error: %v\n", err)
		return errs.ErrInternalServer
	}

	log.Println("[Register] Profile created")

	if err := as.walletRepo.Create(ctx, userID, 0); err != nil {
		log.Printf("[Register] Create wallet error: %v\n", err)
		return errs.ErrInternalServer
	}

	log.Println("[Register] Wallet created")
	log.Println("[Register] SUCCESS")

	return nil
}

func (as *AuthService) SetPin(ctx context.Context, req dto.AddPinRequest) error {
	log.Println("[SetPin] START")

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	hashedPin := hc.Hash(req.Pin)

	if err := as.authRepo.SetPin(ctx, hashedPin, req.UserID); err != nil {
		log.Printf("[SetPin] SetPin error: %v\n", err)
		return errs.ErrInternalServer
	}

	log.Println("[SetPin] SUCCESS")

	return nil
}

func (as *AuthService) CheckPassword(ctx context.Context, pwd string, id int) error {
	log.Printf("[CheckPassword] START userID=%d\n", id)

	oldPwd, err := as.authRepo.GetUserPassword(ctx, id)
	if err != nil {
		log.Printf("[CheckPassword] GetUserPassword error: %v\n", err)
		return errs.ErrInternalServer
	}

	var hc pkg.HashConfig

	if err := hc.Compare(pwd, oldPwd); err != nil {
		log.Printf("[CheckPassword] Invalid password: %v\n", err)
		return errs.ErrInvalidPassword
	}

	log.Println("[CheckPassword] SUCCESS")

	return nil
}

func (as *AuthService) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error {
	log.Printf("[ChangePassword] START userID=%d\n", req.Id)

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	hashedPwd := hc.Hash(req.NewPassword)

	if err := as.authRepo.SetPassword(ctx, hashedPwd, req.Id); err != nil {
		log.Printf("[ChangePassword] SetPassword error: %v\n", err)
		return errs.ErrInternalServer
	}

	log.Println("[ChangePassword] SUCCESS")

	return nil
}

func (as *AuthService) Logout(ctx context.Context, req dto.LogoutRequest) error {
	log.Println("[Logout] START")

	if err := as.authRepo.BlacklistToken(ctx, req.Token, req.ExpiredAt); err != nil {
		log.Printf("[Logout] BlacklistToken error: %v\n", err)
		return err
	}

	log.Println("[Logout] SUCCESS")

	return nil
}

func (as *AuthService) ForgotPassword(ctx context.Context, email string) error {
	log.Printf("[ForgotPassword] START email=%s\n", email)

	clientUrl := os.Getenv("CLIENT_URL")

	user, err := as.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("[ForgotPassword] GetUserByEmail error: %v\n", err)
		return err
	}

	token := utils.GenerateRandomToken()

	if err := as.authRepo.SaveTokenForgotPass(ctx, token, user.Id); err != nil {
		log.Printf("[ForgotPassword] SaveTokenForgotPass error: %v\n", err)
		return err
	}

	resetLink := fmt.Sprintf(
		"%s/reset-password?token=%s",
		clientUrl,
		token,
	)

	log.Printf("[ForgotPassword] Reset link: %s\n", resetLink)
	log.Println("[ForgotPassword] SUCCESS")

	return nil
}

func (as *AuthService) ChangePasswordByReset(ctx context.Context, newPas string, token string) error {
	log.Println("[ChangePasswordByReset] START")

	userID, err := as.authRepo.TokenValidCheck(ctx, token)
	if err != nil {
		log.Printf("[ChangePasswordByReset] TokenValidCheck error: %v\n", err)
		return err
	}

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	hashedPass := hc.Hash(newPas)

	if err := as.authRepo.SetPassword(ctx, hashedPass, userID); err != nil {
		log.Printf("[ChangePasswordByReset] SetPassword error: %v\n", err)
		return err
	}

	if err := as.authRepo.ValidateChangedPassword(ctx, token); err != nil {
		log.Printf("[ChangePasswordByReset] ValidateChangedPassword error: %v\n", err)
		return err
	}

	log.Printf("[ChangePasswordByReset] SUCCESS userID=%d\n", userID)

	return nil
}
