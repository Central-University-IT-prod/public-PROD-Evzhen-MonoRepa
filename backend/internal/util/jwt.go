package util

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var (
	UserSecretKey  = os.Getenv("JWT_USER_SECRET")
	AdminSecretKey = os.Getenv("JWT_ADMIN_SECRET")
)

func GenerateJwt(issuer string, role string) (string, error) {
	var secretKey string

	switch role {
	case RoleUser:
		secretKey = UserSecretKey

	case RoleAdmin:
		secretKey = AdminSecretKey
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	return claims.SignedString([]byte(secretKey))

}

func ParseJwt(jwtToken string, role string) (string, error) {
	var secretKey string

	switch role {
	case RoleUser:
		secretKey = UserSecretKey

	case RoleAdmin:
		secretKey = AdminSecretKey
	}

	token, err := jwt.ParseWithClaims(jwtToken, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	claims := token.Claims.(*jwt.StandardClaims)
	return claims.Issuer, nil
}
