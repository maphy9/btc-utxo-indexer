package middleware

import (
	"net/http"
	"strings"

	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceConfig := helpers.ServiceConfig(r)
		accessTokenKey := serviceConfig.AccessTokenKey

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ape.RenderErr(w, apierrors.NewApiError(
				http.StatusUnauthorized,
				"Authorization header is required",
			))
			return
		}

		accessTokenString := strings.TrimPrefix(authHeader, "Bearer ")
		accessToken, err := helpers.VerifyToken(accessTokenKey, accessTokenString)
		if err != nil {
			ape.RenderErr(w, apierrors.NewApiError(
				http.StatusUnauthorized,
				"The token is invalid",
			))
			return
		}

		userID, err := helpers.GetUserIDFromToken(accessToken)
		if err != nil {
			ape.RenderErr(w, apierrors.NewApiError(
				http.StatusUnauthorized,
				"Could not get the necessary claims from the token",
			))
			return
		}

		ctx := helpers.CtxUserID(userID)(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
