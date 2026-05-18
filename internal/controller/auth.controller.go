package controller

import (
	"backend-golang/internal/dto"
	errs "backend-golang/internal/err"
	"backend-golang/internal/service"
	"backend-golang/pkg"
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
				utils.SendResponse(ctx, http.StatusBadRequest, false, "password is required", nil, err)
				return
			}
			if strings.Contains(err.Error(), "min") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "password must be at least 8 characters", nil, err)
				return
			}
		}
		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid request body", nil, err)
		return
	}
	// jalankan dan kirim ke service
	data, err := ac.authService.Login(ctx.Request.Context(), user)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredential) {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "login failed", nil, err.Error())
			return
		}
		if errors.Is(err, errs.ErrEmailNotFound) {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "login failed", nil, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInternalServer) {
			utils.SendResponse(ctx, http.StatusInternalServerError, false, "login failed", nil, err.Error())
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
				utils.SendResponse(ctx, http.StatusBadRequest, false, "password is required", nil, err)
				return
			}
			if strings.Contains(err.Error(), "min") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "password must be at least 8 characters", nil, err)
				return
			}
		}
		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid request body", nil, err)
		return
	}
	if err := ac.authService.Register(ctx.Request.Context(), user); err != nil {
		if errors.Is(err, errs.ErrExistingEmail) {
			utils.SendResponse(ctx, http.StatusConflict, false, "register failed", nil, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInternalServer) {
			utils.SendResponse(ctx, http.StatusInternalServerError, false, "register failed", nil, err.Error())
			return
		}
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "register failed", nil, err.Error())
		return
	}
	utils.SendResponse(ctx, http.StatusCreated, true, "register success", nil, nil)
}

func (ac *AuthController) SetUserPin(ctx *gin.Context) {
	var user dto.AddPinRequest
	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		if strings.Contains(err.Error(), "required") {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "pin is required", nil, err)
			return
		}
		if strings.Contains(err.Error(), "min") {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "pin must be at least 6 characters", nil, err)
			return
		}
		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid request body", nil, err)
		return
	}
	if err := ac.authService.SetPin(ctx.Request.Context(), user); err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "set pin failed", nil, err.Error())
		return
	}
	utils.SendResponse(ctx, http.StatusCreated, true, "set pin success", nil, nil)

}

func (ac *AuthController) UpdateUserPin(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	var user dto.AddPinRequest
	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		if strings.Contains(err.Error(), "required") {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "pin is required", nil, err)
			return
		}
		if strings.Contains(err.Error(), "min") {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "pin must be at least 6 characters", nil, err)
			return
		}
		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid request body", nil, err.Error())
		return
	}
	user.UserID = claims.Id
	if err := ac.authService.SetPin(ctx.Request.Context(), user); err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "update pin failed", nil, err.Error())
		return
	}
	utils.SendResponse(ctx, http.StatusOK, true, "update pin success", nil, nil)
}

func (ac *AuthController) UpdateUserPassword(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	var user dto.ChangePasswordRequest
	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		if strings.Contains(err.Error(), "required") {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "password is required", nil, err)
			return
		}
		if strings.Contains(err.Error(), "min") {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "password must be at least 8 characters", nil, err)
			return
		}
		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid request body", nil, err.Error())
		return
	}

	if err := ac.authService.CheckPassword(ctx, user.Password, user.Id); err != nil {
		if errors.Is(err, errs.ErrInvalidPassword) {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "update password failed", nil, err.Error())
			return
		}
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "update password failed", nil, err.Error())
		return
	}

	user.Id = claims.Id
	if err := ac.authService.ChangePassword(ctx.Request.Context(), user); err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "update password failed", nil, err.Error())
		return
	}
	utils.SendResponse(ctx, http.StatusOK, true, "update pin success", nil, nil)
}
