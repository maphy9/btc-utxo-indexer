package helpers

import (
	"errors"
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"github.com/maphy9/btc-utxo-indexer/internal/service/requests"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUserRefreshToken(r *http.Request, userID int64, refreshToken string) error {
	ctx := r.Context()
	db := DB(r)
	return db.Users().UpdateRefreshToken(ctx, userID, refreshToken)
}

func RegisterUser(r *http.Request, request requests.RegisterRequest) error {
	ctx := r.Context()
	db := DB(r)

	passwordHash, err := HashPassword(request.Password)
	if err != nil {
		return err
	}

	user := data.User{
		Username:     request.Username,
		PasswordHash: passwordHash,
	}
	_, err = db.Users().Insert(ctx, user)
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(user *data.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

func VerifyUserCredentials(r *http.Request, request requests.LoginRequest) (*data.User, error) {
	ctx := r.Context()
	db := DB(r)

	user, err := db.Users().GetByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("User doesn't exist")
	}
	return user, VerifyPassword(user, request.Password)
}

func GetUserRefreshToken(r *http.Request, userID int64) (string, error) {
	ctx := r.Context()
	db := DB(r)

	user, err := db.Users().GetByUserID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("User doesn't exist")
	}
	return user.RefreshToken, nil
}
