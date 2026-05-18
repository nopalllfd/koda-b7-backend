package service

import (
	"backend-golang/internal/dto"
	errs "backend-golang/internal/err"
	"backend-golang/internal/repository"
	"context"
	"errors"
	"log"
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
		if errors.Is(err, errs.ErrProfileNotFound) {
			return dto.Profiles{}, err
		}
		if errors.Is(err, errs.ErrInternalServer) {
			return dto.Profiles{}, err
		}
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
		if errors.Is(err, errs.ErrPhoneAlreadyUsed) || errors.Is(err, errs.ErrInvalidInput) || errors.Is(err, errs.ErrInternalServer) {
			log.Printf("[ERROR] database failure during registration: %v", err)
			return err
		}
		return errs.ErrInternalServer
	}

	if result == 0 {
		return errs.ErrProfileNotFound
	}
	return nil
}
