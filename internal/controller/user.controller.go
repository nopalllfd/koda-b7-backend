package controller

import (
	"backend-golang/internal/dto"
	errs "backend-golang/internal/err"
	"backend-golang/internal/service"
	"backend-golang/pkg"
	"backend-golang/pkg/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Get User Profile
//
//	@Summary		Get user profile
//	@Description	get authenticated user profile
//	@Tags			user
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	dto.ProfileSwaggerResponse
//	@Failure		401	{object}	dto.ErrorSwaggerResponse
//	@Failure		404	{object}	dto.ErrorSwaggerResponse
//	@Failure		500	{object}	dto.ErrorSwaggerResponse
//	@Router			/user/profile [get]
func (uc *UserController) GetProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	user, err := uc.userService.GetUserProfile(ctx.Request.Context(), claims.Id)
	if err != nil {
		log.Println(err.Error())

		if errors.Is(err, errs.ErrProfileNotFound) {
			utils.SendResponse(ctx, http.StatusNotFound, false, "get profile failed", nil, err.Error())
			return
		}

		utils.SendResponse(ctx, http.StatusInternalServerError, false, "get profile failed", nil, err.Error())
		return
	}

	utils.SendResponse(ctx, http.StatusOK, true, "ok", user, nil)
}

// Edit User Profile
//
//	@Summary		Edit user profile
//	@Description	update authenticated user profile
//	@Tags			user
//	@Security		ApiKeyAuth
//	@Accept			mpfd
//	@Produce		json
//	@Param			fullname	formData	string	false	"full name"
//	@Param			phone		formData	string	false	"phone number"
//	@Param			photo		formData	file	false	"profile photo"
//	@Success		200			{object}	dto.ProfileSwaggerResponse
//	@Failure		400			{object}	dto.ErrorSwaggerResponse
//	@Failure		401			{object}	dto.ErrorSwaggerResponse
//	@Failure		404			{object}	dto.ErrorSwaggerResponse
//	@Failure		409			{object}	dto.ErrorSwaggerResponse
//	@Failure		500			{object}	dto.ErrorSwaggerResponse
//	@Router			/user/profile [put]
func (uc *UserController) EditProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var profile dto.ProfileUpdateRequest

	if err := ctx.ShouldBindWith(&profile, binding.FormMultipart); err != nil {
		utils.SendResponse(
			ctx,
			http.StatusBadRequest,
			false,
			"invalid request body",
			nil,
			err.Error(),
		)
		return
	}

	if profile.Photo != nil {

		ext := path.Ext(profile.Photo.Filename)

		filename := fmt.Sprintf(
			"profile_%d%s",
			time.Now().UnixNano(),
			ext,
		)

		dst := filepath.Join("public", "img", filename)

		if err := ctx.SaveUploadedFile(profile.Photo, dst); err != nil {
			utils.SendResponse(
				ctx,
				http.StatusInternalServerError,
				false,
				"edit profile failed",
				nil,
				err.Error(),
			)
			return
		}

		profile.PhotoPath = dst
	}

	if err := uc.userService.EditProfile(
		ctx.Request.Context(),
		claims.Id,
		profile,
	); err != nil {

		if errors.Is(err, errs.ErrPhoneAlreadyUsed) {
			utils.SendResponse(ctx, http.StatusConflict, false, "edit profile failed", nil, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProfileNotFound) {
			utils.SendResponse(ctx, http.StatusNotFound, false, "edit profile failed", nil, err.Error())
			return
		}

		if errors.Is(err, errs.ErrInvalidInput) {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "edit profile failed", nil, err.Error())
			return
		}

		utils.SendResponse(ctx, http.StatusInternalServerError, false, "edit profile failed", nil, err.Error())
		return
	}

	utils.SendResponse(
		ctx,
		http.StatusOK,
		true,
		"ok",
		gin.H{
			"fullname": profile.FullName,
			"phone":    profile.Phone,
			"photo":    profile.PhotoPath,
		},
		nil,
	)
}
