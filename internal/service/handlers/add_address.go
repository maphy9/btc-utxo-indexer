package handlers

import (
	"errors"
	"net/http"

	"github.com/lib/pq"
	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
)

func AddAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := helpers.Log(r)
	db := helpers.DB(r)
	manager := helpers.Manager(r)
	userID := helpers.UserID(r)

	request, err := requests.NewAddAddressRequest(r)
	if err != nil {
		ape.RenderErr(w, apierrors.BadRequest())
		return
	}

	err = helpers.AddAddress(ctx, db, manager, userID, request.Address)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ape.RenderErr(w, apierrors.NewApiError(
				http.StatusConflict,
				"This address is already being tracked",
			))
		} else {
			logger.WithError(err).Debug("Failed to add an address for tracking")
			ape.RenderErr(w, apierrors.NewApiError(
				http.StatusInternalServerError,
				"Failed to add an address",
			))
		}
		return
	}
}
