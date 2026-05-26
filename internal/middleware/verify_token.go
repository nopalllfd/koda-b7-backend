package middleware

import (
	"backend-golang/internal/dto"
	"backend-golang/pkg"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func VerifyMiddleware(db *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
				Error:   "invalid token format",
			})
			return
		}

		token := splittedBearer[1]

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

		// CEK TOKEN BLACKLIST KE DATABASE
		var isBlacklisted bool
		sql := `SELECT EXISTS(SELECT 1 FROM token_blacklists WHERE token = $1)`
		err := db.QueryRow(ctx.Request.Context(), sql, token).Scan(&isBlacklisted)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
				Message: "Error",
				Success: false,
				Error:   "failed to verify token status",
			})
			return
		}

		if isBlacklisted {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: "Unauthorized Access, Please Login",
				Success: false,
				Error:   "token has been revoked",
			})
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
