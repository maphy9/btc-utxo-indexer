package helpers

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateJWTTokens(r *http.Request, userID int64) (string, string, error) {
	serviceConfig := ServiceConfig(r)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(serviceConfig.TokenExpireTime).Unix(),
			"type":    "access",
		},
	)
	tokenString, err := token.SignedString([]byte(serviceConfig.TokenKey))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(serviceConfig.RefreshTokenExpireTime).Unix(),
			"type":    "refresh",
		},
	)
	refreshTokenString, err := refreshToken.SignedString([]byte(serviceConfig.RefreshTokenKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, err
}

func VerifyToken(r *http.Request, tokenString string) (*jwt.Token, error) {
	serviceConfig := ServiceConfig(r)
	tokenKey := serviceConfig.TokenKey

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return []byte(tokenKey), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return token, nil
}

func GetUserIDFromToken(token *jwt.Token) (int64, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("The token doesn't have the necessary claims")
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("The token doesn't have the necessary claims")
	}
	return int64(userID), nil
}
