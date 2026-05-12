package main

import (
	// untuk logging
	// informasi dari protocol (status code)

	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin" // framework gin gonic
	"github.com/gin-gonic/gin/binding"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   any    `json:"error"`
}

func SendResponse(c *gin.Context, code int, success bool, message string, data any, err any) {
	c.JSON(code, Response{
		Success: success,
		Message: message,
		Data:    data,
		Error:   err,
	})

	log.Printf(
		"[%s] %s -> %d | success=%v | message=%s",
		c.Request.Method,
		c.Request.URL.Path,
		code,
		success,
		message,
	)
}

func ValidateEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regex.MatchString(email)
}

func main() {
	// init
	app := gin.Default()

	// route
	//login
	app.POST("/auth/login", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
			SendResponse(ctx, http.StatusInternalServerError, false, "Internal Server Error", nil, err)
			return
		}
		email := user.Email
		password := user.Password

		// validation

		// email format validation
		emailValidate := ValidateEmail(email)
		if !emailValidate {
			err := gin.H{
				"email": "invalid email format",
			}
			SendResponse(ctx, http.StatusBadRequest, false, "Bad Request", nil, err)
			return
		}

		// password length validation
		if len(password) < 7 {
			err := gin.H{
				"password": "password must be at least 7 characters",
			}
			SendResponse(ctx, http.StatusBadRequest, false, "Bad Request", nil, err)
			return
		}

		// credentials validation
		if email != "nopal@gmail.com" || password != "1234567" {
			SendResponse(ctx, http.StatusUnauthorized, false, "Invalid email or password", nil, nil)
			return
		}

		data := gin.H{
			"user": gin.H{
				"email": email,
			},
		}
		SendResponse(ctx, http.StatusOK, true, "Login Success", data, nil)

	})
	//register
	app.POST("/auth/register", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
			SendResponse(ctx, http.StatusInternalServerError, false, "Internal Server Error", nil, err)
			return
		}
		email := user.Email
		password := user.Password

		// validation
		// email format validation
		emailValidate := ValidateEmail(email)
		if !emailValidate {
			err := gin.H{
				"email": "invalid email format",
			}
			SendResponse(ctx, http.StatusBadRequest, false, "Bad Request", nil, err)
			return
		}
		// password length validation
		if len(password) < 7 {
			err := gin.H{
				"password": "password must be at least 7 characters",
			}
			SendResponse(ctx, http.StatusBadRequest, false, "Bad Request", nil, err)
			return
		}
		//existing email validation
		if email == "nopal@gmail.com" {
			SendResponse(ctx, http.StatusBadRequest, false, "Email already exists", nil, nil)
			return
		}

		data := gin.H{
			"user": gin.H{
				"email": email,
			},
		}

		SendResponse(ctx, http.StatusCreated, true, "Register Success", data, nil)

	})
	app.Run("localhost:5000")
}
