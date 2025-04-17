package crypto

import (
	"errors"
	"server/internal/constant"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID       uint
	Username string
	jwt.RegisteredClaims
}

var secret = []byte(constant.JWTConfig.Secret)

func GenerateJWT(id uint, username string, expiresAt time.Duration, subject string) (string, error) {
	claims := JWTClaims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateJWTToKanboard(id uint, username string) (string, error) {
	return GenerateJWT(id, username, constant.JWTConfig.KanboardTokenExpiration, constant.KANBOARD_SUBJECT)
}

func GenerateJWTToAdmin(id uint, username string) (string, error) {
	return GenerateJWT(id, username, constant.JWTConfig.AdminTokenExpiration, constant.ADMIN_SUBJECT)
}

func ParseToken(tokenStr string) (*JWTClaims, error) {
	var claims JWTClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err == nil && !token.Valid {
		err = errors.New("invalid token")
	}

	return &claims, err
}

func IsTokenValid(tokenStr string) bool {
	_, err := ParseToken(tokenStr)
	return err == nil
}
