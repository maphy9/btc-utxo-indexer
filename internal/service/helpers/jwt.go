package helpers

import (
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
	tokenString, err := token.SignedString(serviceConfig.TokenKey)
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
	refreshTokenString, err := refreshToken.SignedString(serviceConfig.RefreshTokenKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, err
}

func UpdateUserRefreshToken(r *http.Request, userID int64, refreshToken string) error {
	ctx := r.Context()
	db := DB(r)
	return db.Users().UpdateRefreshToken(ctx, userID, refreshToken)
}