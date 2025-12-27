package handlers

import (
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/requests"
	"github.com/maphy9/btc-utxo-indexer/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
)

func Login(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)
	request, err := requests.NewLoginRequest(r)
	if err != nil {
		logger.WithError(err).Debug("bad request")
		ape.RenderErr(w, apierrors.BadRequest())
		return
	}

	user, err := helpers.VerifyUserCredentials(r, request.Username, request.Password)
	if err != nil {
		logger.WithError(err).Debug("invalid user credentials")
		ape.RenderErr(w, apierrors.NewApiError(http.StatusForbidden, "Invalid user credentials"))
		return
	}

	accessToken, refreshToken, err := helpers.GenerateJWTTokens(r, user.ID)
	if err != nil {
		logger.WithError(err).Error("failed to generate tokens")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError, "Failed to generate tokens",
		))
		return
	}

	err = helpers.UpdateUserRefreshToken(r, user.ID, refreshToken)
	if err != nil {
		logger.WithError(err).Error("failed to update the refresh token")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError, "Failed to update the refresh token",
		))
		return
	}

	ape.Render(w, responses.NewTokenResponse(accessToken, refreshToken))
}
