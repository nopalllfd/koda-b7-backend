package service

import (
	"backend-golang/internal/dto"
	"backend-golang/internal/repository"
	"context"
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
	log.Println(&data.FullName)
	if err := us.userRepo.Edit(ctx, &data.FullName, &data.Photo, &data.Phone, id); err != nil {
		return err
	}
	return nil
}
