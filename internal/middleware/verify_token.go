package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/nopalllfd/koda-b7-backend/internal/dto"
	"github.com/nopalllfd/koda-b7-backend/pkg"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func VerifyMiddleware(
	rc *redis.Client,
) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		bearerToken := ctx.GetHeader("Authorization")

		if bearerToken == "" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				dto.Response{
					Message: "Unauthorized Access, Please Login",
					Success: false,
					Error:   "unauthorized access, please login",
				},
			)
			return
		}

		splittedBearer := strings.Split(
			bearerToken,
			" ",
		)

		if len(splittedBearer) != 2 ||
			splittedBearer[0] != "Bearer" {

			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				dto.Response{
					Message: "Unauthorized Access, Please Login",
					Success: false,
					Error:   "invalid token format",
				},
			)
			return
		}

		token := splittedBearer[1]

		var claims pkg.Claims

		if err := claims.VerifyJWT(token); err != nil {

			if errors.Is(err, jwt.ErrTokenInvalidIssuer) ||
				errors.Is(err, jwt.ErrTokenExpired) {

				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					dto.Response{
						Message: "Unauthorized Access, Please Login",
						Success: false,
						Error:   err.Error(),
					},
				)
				return
			}

			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				dto.Response{
					Message: "Error",
					Success: false,
					Error:   "Internal Server Error",
				},
			)
			return
		}

		// CEK TOKEN BLACKLIST KE REDIS
		result, err := rc.Exists(
			ctx.Request.Context(),
			"bl:"+token,
		).Result()

		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				dto.Response{
					Message: "Error",
					Success: false,
					Error:   "failed to verify token status",
				},
			)
			return
		}

		if result > 0 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				dto.Response{
					Message: "Unauthorized Access, Please Login",
					Success: false,
					Error:   "token has been revoked",
				},
			)
			return
		}

		ctx.Set("claims", claims)

		ctx.Next()
	}
}
