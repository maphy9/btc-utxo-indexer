package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
)

func GetUtxos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := helpers.Log(r)
	db := helpers.DB(r)
	userID := helpers.UserID(r)
	address := chi.URLParam(r, "address")

	found, err := helpers.CheckAddress(ctx, db, userID, address)
	if err != nil {
		logger.WithError(err).Error("failed to check the address")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError,
			"Failed to check the address",
		))
		return
	}
	if !found {
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusNotFound,
			"Address not found",
		))
		return
	}

	utxos, err := helpers.GetUtxos(ctx, db, address)
	if err != nil {
		logger.WithError(err).Error("failed to get utxos")
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError,
			"Failed to get utxos",
		))
		return
	}

	ape.Render(w, responses.NewGetUtxosResponse(utxos))
}
