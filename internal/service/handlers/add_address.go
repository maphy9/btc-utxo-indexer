package handlers

import (
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/requests"
	"github.com/maphy9/btc-utxo-indexer/internal/util"
	"gitlab.com/distributed_lab/ape"
)

func AddAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := helpers.Log(r)
	db := helpers.DB(r)
	userID := helpers.UserID(r)
	manager := helpers.Manager(r)

	request, err := requests.NewAddAddressRequest(r)
	if err != nil {
		ape.RenderErr(w, apierrors.BadRequest())
		return
	}

	err = helpers.AddAddress(ctx, db, userID, request.Address)
	if err != nil {
		if util.IsUniqueViolation(err) {
			ape.RenderErr(w, apierrors.NewApiError(
				http.StatusConflict,
				"This address is already being tracked",
			))
		} else {
			logger.WithError(err).Debug("failed to add an address for tracking")
			ape.RenderErr(w, apierrors.NewApiError(
				http.StatusInternalServerError,
				"Failed to add an address",
			))
		}
		return
	}

	if err := manager.SubscribeAddress(request.Address); err != nil {
		logger.WithError(err).Error("failed to subscribe to an address")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError,
			"Failed to subscribe to an address",
		))
		return
	}
}
