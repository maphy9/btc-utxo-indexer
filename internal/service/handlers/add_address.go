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
	logger := helpers.Log(r)
	request, err := requests.NewAddAddressRequest(r)
	if err != nil {
		ape.RenderErr(w, apierrors.BadRequest())
		return
	}

	err = helpers.AddAddress(r, request.Address)
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
