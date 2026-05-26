package controller

import (
	"backend-golang/internal/dto"
	errs "backend-golang/internal/err"
	"backend-golang/internal/service"
	"backend-golang/pkg"
	"backend-golang/pkg/utils"
	"errors"
	"log"
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

// User Login
//
//	@Summary		Login user
//	@Description	login user using email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.LoginRequest	true	"login payload"
//	@Success		200		{object}	dto.LoginSwaggerResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/auth/login [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var user dto.LoginRequest

	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {

		if strings.Contains(err.Error(), "Email") {

			if strings.Contains(err.Error(), "required") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "email is required", nil, err.Error())
				return
			}

			if strings.Contains(err.Error(), "email") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid email format", nil, err.Error())
				return
			}
		}

		if strings.Contains(err.Error(), "Password") {

			if strings.Contains(err.Error(), "required") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "password is required", nil, err.Error())
				return
			}

			if strings.Contains(err.Error(), "min") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "password must be at least 8 characters", nil, err)
				return
			}
		}

		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid request body", nil, err.Error())
		return
	}

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

	utils.SendResponse(ctx, http.StatusOK, true, "login success", data, nil)
}

// User Register
//
//	@Summary		Register user
//	@Description	register new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.RegisterRequest	true	"register payload"
//	@Success		201		{object}	dto.RegisterSwaggerResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		409		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/auth/register [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var user dto.RegisterRequest

	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {

		if strings.Contains(err.Error(), "Email") {

			if strings.Contains(err.Error(), "required") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "email is required", nil, err.Error())
				return
			}

			if strings.Contains(err.Error(), "email") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid email format", nil, err.Error())
				return
			}
		}

		if strings.Contains(err.Error(), "Password") {

			if strings.Contains(err.Error(), "required") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "password is required", nil, err.Error())
				return
			}

			if strings.Contains(err.Error(), "min") {
				utils.SendResponse(ctx, http.StatusBadRequest, false, "password must be at least 8 characters", nil, err)
				return
			}
		}

		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid request body", nil, err.Error())
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

// Set User PIN
//
//	@Summary		Set user PIN
//	@Description	create PIN for user
//	@Tags			auth
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UserPIN	true	"pin payload"
//	@Success		201		{object}	dto.RegisterSwaggerResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/auth/pin [post]
func (ac *AuthController) SetUserPin(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var user dto.AddPinRequest

	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {

		if strings.Contains(err.Error(), "required") {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "pin is required", nil, nil)
			return
		}

		if strings.Contains(err.Error(), "min") {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "pin must be at least 6 characters", nil, nil)
			return
		}

		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid request body", nil, err.Error())
		return
	}
	user.UserID = claims.Id

	if err := ac.authService.SetPin(ctx.Request.Context(), user); err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "set pin failed", nil, err.Error())
		return
	}

	utils.SendResponse(ctx, http.StatusCreated, true, "set pin success", nil, nil)
}

// Update User PIN
//
//	@Summary		Update user PIN
//	@Description	update existing user PIN
//	@Tags			auth
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UserPIN	true	"pin payload"
//	@Success		200		{object}	dto.RegisterSwaggerResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		401		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/auth/pin [patch]
func (ac *AuthController) UpdateUserPin(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var user dto.AddPinRequest

	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {

		if strings.Contains(err.Error(), "required") {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "pin is required", nil, err.Error())
			return
		}

		if strings.Contains(err.Error(), "min") {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "pin must be at least 6 characters", nil, err.Error())
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

// Update User Password
//
//	@Summary		Update user password
//	@Description	change current user password
//	@Tags			auth
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.ChangePasswordRequest	true	"change password payload"
//	@Success		200		{object}	dto.RegisterSwaggerResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		401		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/auth/password [patch]
func (ac *AuthController) UpdateUserPassword(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var user dto.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	user.Id = claims.Id

	log.Println("cek pw controller")

	if err := ac.authService.CheckPassword(
		ctx.Request.Context(),
		user.OldPassword,
		user.Id,
	); err != nil {

		if errors.Is(err, errs.ErrInvalidPassword) {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "old password invalid", nil, err.Error())
			return
		}

		utils.SendResponse(ctx, http.StatusInternalServerError, false, "update password failed", nil, err.Error())
		return
	}

	log.Println("Abis cek pw")

	if err := ac.authService.ChangePassword(ctx.Request.Context(), user); err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "update password failed", nil, err.Error())
		return
	}

	utils.SendResponse(ctx, http.StatusOK, true, "update password success", nil, nil)
}

// Edit User Profile
//
//	@Summary		Delete user token
//	@Description	for logout
//	@Tags			auth
//	@Security		ApiKeyAuth
//	@Success		200			{object}	dto.LogoutSwaggerResponse
//	@Failure		401			{object}	dto.ErrorSwaggerResponse
//	@Failure		404			{object}	dto.ErrorSwaggerResponse
//	@Failure		409			{object}	dto.ErrorSwaggerResponse
//	@Failure		500			{object}	dto.ErrorSwaggerResponse
//	@Router			/auth/logout [delete]
func (ac *AuthController) Logout(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	authHeader := ctx.GetHeader("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	data := dto.LogoutRequest{
		Token:     tokenString,
		ExpiredAt: claims.ExpiresAt.Time,
	}
	ac.authService.Logout(ctx.Request.Context(), data)
}

// Forgot Password
//
//	@Summary		Forgot password
//	@Description	send reset password link to email
//	@Tags			auth
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.ForgotPasswordRequest	true	"forgot password payload"
//	@Success		200		{object}	dto.RegisterSwaggerResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/auth/forgot-password [post]
func (ac *AuthController) ForgotPassword(ctx *gin.Context) {

	var req dto.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		if strings.Contains(err.Error(), "Email") {

			if strings.Contains(err.Error(), "required") {
				utils.SendResponse(
					ctx,
					http.StatusBadRequest,
					false,
					"email is required",
					nil,
					err.Error(),
				)
				return
			}

			if strings.Contains(err.Error(), "email") {
				utils.SendResponse(
					ctx,
					http.StatusBadRequest,
					false,
					"invalid email format",
					nil,
					err.Error(),
				)
				return
			}
		}

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

	if err := ac.authService.ForgotPassword(
		ctx.Request.Context(),
		req.Email,
	); err != nil {

		utils.SendResponse(
			ctx,
			http.StatusInternalServerError,
			false,
			"forgot password failed",
			nil,
			err.Error(),
		)
		return
	}

	utils.SendResponse(
		ctx,
		http.StatusOK,
		true,
		"reset password link sent",
		nil,
		nil,
	)
}

// Reset Password
//
//	@Summary		Reset password
//	@Description	reset password using reset token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			token	query		string	true	"reset token"
//	@Param			body	body		dto.ResetPasswordRequest	true	"reset password payload"
//	@Success		200		{object}	dto.RegisterSwaggerResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/auth/reset-password [patch]
func (ac *AuthController) ResetPassword(ctx *gin.Context) {

	var req dto.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	if err := ac.authService.ChangePasswordByReset(
		ctx.Request.Context(),
		req.NewPassword,
		req.Token,
	); err != nil {

		utils.SendResponse(
			ctx,
			http.StatusBadRequest,
			false,
			"reset password failed",
			nil,
			err.Error(),
		)
		return
	}

	utils.SendResponse(
		ctx,
		http.StatusOK,
		true,
		"reset password success",
		nil,
		nil,
	)
}
