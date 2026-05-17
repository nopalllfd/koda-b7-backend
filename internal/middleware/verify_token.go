package middleware

import (
	"backend-golang/internal/dto"
	"backend-golang/pkg"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyMiddleware(ctx *gin.Context) {
	// mengambil token dari request payload (header)
	// header Authorization
	// Bearer token => token diawali dengan kata Bearer
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Message: "Unauthorized Access, Please Login",
			Success: false,
			Error:   "unauthorized access, please login",
		})
		return
	}
	splittedBearer := strings.Split(bearerToken, " ")

	if len(splittedBearer) != 2 || splittedBearer[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Message: "Unauthorized Access, Please Login",
			Success: false,
			Error:   "invalid token",
		})
		return
	}
	if len(splittedBearer) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Message: "Unauthorized Access, Please Login",
			Success: false,
			Error:   "invalid token",
		})
		return
	}
	token := splittedBearer[1]

	// verifikasi token
	var claims pkg.Claims
	if err := claims.VerifyJWT(token); err != nil {
		if errors.Is(err, jwt.ErrTokenInvalidIssuer) || errors.Is(err, jwt.ErrTokenExpired) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: "Unauthorized Access, Please Login",
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
			Message: "Error",
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}

	ctx.Set("claims", claims)
	ctx.Next()

}
