package service

import (
	"backend-golang/internal/dto"
	errs "backend-golang/internal/err"
	"backend-golang/internal/repository"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (us *UserService) GetUserProfile(ctx context.Context, id int) (dto.Profiles, error) {
	result, err := us.userRepo.GetProfile(ctx, id)
	if err != nil {
		return dto.Profiles{}, err
	}

	return dto.Profiles{
		User_id:   result.User_id,
		FullName:  result.FullName,
		Photo:     result.Photo,
		Phone:     result.Phone,
		CreatedAt: result.Created_at,
		UpdatedAt: &result.Updated_at,
	}, nil
}

func (us *UserService) EditProfile(ctx context.Context, id int, data dto.ProfileUpdateRequest) error {
	log.Println(data.FullName)
	result, err := us.userRepo.Edit(ctx, &data.FullName, &data.Photo, &data.Phone, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // Kode Postgres untuk Unique Violation
				return errs.ErrPhoneAlreadyUsed
			case "22001": // Kode Postgres untuk String Data Right Truncation (input terlalu panjang)
				return errs.ErrInvalidInput
			}
		}
		// Error database lainnya (koneksi putus, dll)
		return errs.ErrInternalServer
	}

	if result == 0 {
		return errs.ErrProfileNotFound
	}
	return nil
}
