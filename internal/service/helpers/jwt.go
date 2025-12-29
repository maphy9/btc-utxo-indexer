package helpers

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateJWTTokens(r *http.Request, userID int64) (string, string, error) {
	serviceConfig := ServiceConfig(r)

	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(serviceConfig.AccessTokenExpireTime).Unix(),
			"type":    "access",
		},
	)
	tokenString, err := accessToken.SignedString([]byte(serviceConfig.AccessTokenKey))
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

func verifyToken(tokenString, tokenKey string) (*jwt.Token, error) {
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
		return nil, errors.New("invalid token")
	}
	return token, nil
}

func VerifyAccessToken(r *http.Request, accessTokenString string) (*jwt.Token, error) {
	serviceConfig := ServiceConfig(r)
	accessTokenKey := serviceConfig.AccessTokenKey
	return verifyToken(accessTokenString, accessTokenKey)
}

func VerifyRefreshToken(r *http.Request, refreshTokenString string) (*jwt.Token, error) {
	serviceConfig := ServiceConfig(r)
	refreshTokenKey := serviceConfig.RefreshTokenKey
	return verifyToken(refreshTokenString, refreshTokenKey)
}

func GetUserIDFromToken(token *jwt.Token) (int64, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("the token doesn't have the necessary claims")
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("the token doesn't have the necessary claims")
	}
	return int64(userID), nil
}
