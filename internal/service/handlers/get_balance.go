package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := helpers.Log(r)
	db := helpers.DB(r)
	address := chi.URLParam(r, "address")

	balance, err := db.Transactions().GetAddressBalance(ctx, address)
	if err != nil {
		logger.WithError(err).Error("failed to get address balance")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError,
			"Failed to get address balance",
		))
		return
	}

	ape.Render(w, responses.NewGetBalanceResponse(balance))
}
