package handlers

import (
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
)

func GetAddresses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := helpers.Log(r)
	db := helpers.DB(r)
	userID := helpers.UserID(r)

	addresses, err := helpers.GetAddresses(ctx, db, userID)
	if err != nil {
		logger.WithError(err).Debug("Failed to get tracked addresses")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError,
			"Failed to get tracked addresses",
		))
		return
	}

	ape.Render(w, responses.NewGetAddressesResponse(addresses))
}
