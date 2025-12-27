package handlers

import (
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/requests"
	"github.com/maphy9/btc-utxo-indexer/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)
	request, err := requests.NewRefreshRequest(r)
	if err != nil {
		logger.WithError(err).Debug("bad request")
		ape.RenderErr(w, apierrors.BadRequest())
		return
	}

	refreshToken, err := helpers.VerifyRefreshToken(r, request.RefreshToken)
	if err != nil {
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusUnauthorized,
			"The token is invalid",
		))
		return
	}

	userID, err := helpers.GetUserIDFromToken(refreshToken)
	if err != nil {
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusUnauthorized,
			"Could not get the necessary claims from the token",
		))
		return
	}

	savedRefreshToken, err := helpers.GetUserRefreshToken(r, userID)
	if err != nil {
		logger.WithError(err).Debug("Could not retrieve saved refresh token")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError,
			"Could not retrieve saved refresh token",
		))
		return
	}
	if request.RefreshToken != savedRefreshToken {
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusUnauthorized,
			"Token revoked",
		))
		return
	}

	newAccessToken, newRefreshToken, err := helpers.GenerateJWTTokens(r, userID)
	if err != nil {
		logger.WithError(err).Error("failed to generate tokens")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError, "Failed to generate tokens",
		))
		return
	}

	err = helpers.UpdateUserRefreshToken(r, userID, newRefreshToken)
	if err != nil {
		logger.WithError(err).Error("failed to update the refresh token")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError, "Failed to update the refresh token",
		))
		return
	}

	ape.Render(w, responses.NewTokenResponse(newAccessToken, newRefreshToken))
}
