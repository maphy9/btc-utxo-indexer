package helpers

import (
	"errors"
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"github.com/maphy9/btc-utxo-indexer/internal/service/requests"
	"golang.org/x/crypto/bcrypt"
)

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