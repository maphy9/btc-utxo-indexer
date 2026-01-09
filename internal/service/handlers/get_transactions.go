package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
)

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := helpers.Log(r)
	db := helpers.DB(r)
	address := chi.URLParam(r, "address")

	transactions, err := db.Transactions().GetAddressTransactions(ctx, address)
	if err != nil {
		logger.WithError(err).Error("failed to get address transactions")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError,
			"Failed to get address transactions",
		))
		return
	}

	ape.Render(w, transactions)
}
