package service

import (
	"context"
	"errors"
	"log"

	"github.com/nopalllfd/koda-b7-backend/internal/dto"
	errs "github.com/nopalllfd/koda-b7-backend/internal/err"
	"github.com/nopalllfd/koda-b7-backend/internal/repository"
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

func (us *UserService) EditProfile(
	ctx context.Context,
	userID int,
	req dto.ProfileUpdateRequest,
) error {

	if req.Phone != "" {

		exists, err := us.userRepo.FindByPhone(
			ctx,
			req.Phone,
			userID,
		)

		if err != nil {
			return err
		}

		if exists {
			return errs.ErrPhoneAlreadyUsed
		}
	}

	var fullName *string
	var phone *string
	var photo *string

	if req.FullName != "" {
		fullName = &req.FullName
	}

	if req.Phone != "" {
		phone = &req.Phone
	}

	if req.PhotoPath != "" {
		photo = &req.PhotoPath
	}

	rowsAffected, err := us.userRepo.Edit(
		ctx,
		fullName,
		photo,
		phone,
		userID,
	)

	if err != nil {
		log.Println("SERVICE ERROR:", err)
		return err
	}

	if rowsAffected == 0 {
		return errs.ErrProfileNotFound
	}

	return nil
}
