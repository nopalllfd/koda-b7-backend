package controller

import (
	"backend-golang/internal/dto"
	"backend-golang/internal/service"
	"backend-golang/pkg"
	"backend-golang/pkg/utils"
	"log"
	"net/http"

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

func (uc *UserController) GetProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	user, err := uc.userService.GetUserProfile(ctx.Request.Context(), claims.Id)
	if err != nil {
		log.Println(err.Error())
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "internal server error", nil, err.Error())
		return
	}
	utils.SendResponse(ctx, http.StatusOK, true, "ok", user, nil)
}

func (uc *UserController) EditProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	var profile dto.ProfileUpdateRequest
	if err := ctx.ShouldBindBodyWith(&profile, binding.JSON); err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "Internal Server Error", nil, err.Error())
		return
	}

	if err := uc.userService.EditProfile(ctx.Request.Context(), claims.Id, profile); err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "Internal Server Error", nil, err.Error())
		return
	}
	utils.SendResponse(ctx, http.StatusOK, true, "ok", profile, nil)
}
