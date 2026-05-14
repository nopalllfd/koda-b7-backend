package controller

import (
	"backend-golang/internal/dto"
	"backend-golang/internal/service"
	"backend-golang/pkg/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var user dto.LoginRequest
	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		// apakah error validasi ?
		if strings.Contains(err.Error(), "Email") {
			if strings.Contains(err.Error(), "required") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "email is required", nil, err)
				return
			}

			if strings.Contains(err.Error(), "email") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid email format", nil, err)
				return
			}
		}
		if strings.Contains(err.Error(), "Password") {
			if strings.Contains(err.Error(), "required") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "email is required", nil, err)
				return
			}
			if strings.Contains(err.Error(), "min") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "password must be at least 8 characters", nil, err)
				return
			}
		}
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "Internal Server Error", nil, err)
		return
	}
	// jalankan dan kirim ke service
	data, err := ac.authService.Login(ctx.Request.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredential) {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "login failed", nil, err.Error())
			return
		}
		if errors.Is(err, service.ErrEmailNotFound) {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "login failed", nil, err.Error())
			return
		}
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "login failed", nil, err.Error())
		return
	}
	utils.SendResponse(ctx, http.StatusOK, true, "login Success", data, nil)
}

func (ac *AuthController) Register(ctx *gin.Context) {
	var user dto.RegisterRequest
	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		// apakah error validasi ?
		if strings.Contains(err.Error(), "Email") {
			if strings.Contains(err.Error(), "required") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "email is required", nil, err)
				return
			}

			if strings.Contains(err.Error(), "email") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid email format", nil, err)
				return
			}
		}
		if strings.Contains(err.Error(), "Password") {
			if strings.Contains(err.Error(), "required") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "email is required", nil, err)
				return
			}
			if strings.Contains(err.Error(), "min") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "password must be at least 8 characters", nil, err)
				return
			}
		}
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "Internal Server Error", nil, err)
		return
	}
	if err := ac.authService.Register(ctx.Request.Context(), user); err != nil {
		if errors.Is(err, service.ErrExistingEmail) {
			utils.SendResponse(ctx, http.StatusConflict, false, "register Failed", nil, err.Error())
			return
		}
		if errors.Is(err, service.ErrInternalServer) {
			utils.SendResponse(ctx, http.StatusInternalServerError, false, "register Failed", nil, err.Error())
			return
		}
	}
	utils.SendResponse(ctx, http.StatusCreated, true, "register Success", nil, nil)
}
