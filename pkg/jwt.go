package pkg

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Id    int
	Email string
	jwt.RegisteredClaims
}

func NewClaims(id int, email string) *Claims {
	return &Claims{
		Id:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
	}
}

func (c *Claims) GenJWT() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("missing secret key")
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return unsignedToken.SignedString([]byte(jwtSecret))
}

func (c *Claims) VerifyJWT(token string) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("missing secret key")
	}
	log.Println("Token", token)
	jwtToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return jwt.ErrTokenExpired
	}

	issuer, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return err
	}

	if issuer != os.Getenv("JWT_ISSUER") {
		return jwt.ErrTokenInvalidIssuer
	}

	return nil
}
