package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JwtKey = []byte("my_secret_key")
var RefreshTokenKey = []byte("my_refresh_token_key")

type Claims struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(email string, userId int) (string, error) {
	expirationTime := time.Now().Add(3 * time.Minute)

	claims := &Claims{
		UserId: userId,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(email string) (string, error) {
	expirationTime := time.Now().Add(4 * time.Minute)

	claims := &RefreshClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(RefreshTokenKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return RefreshTokenKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
