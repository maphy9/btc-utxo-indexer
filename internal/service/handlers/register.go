package handlers

import (
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/data/pg"
	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
)

func Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := helpers.Log(r)
	db := helpers.DB(r)

	request, err := requests.NewRegisterRequest(r)
	if err != nil {
		ape.RenderErr(w, apierrors.BadRequest())
		return
	}

	err = helpers.RegisterUser(ctx, db, request.Username, request.Password)
	if err != nil {
		if pg.IsUniqueViolation(err) {
			ape.RenderErr(w, apierrors.NewApiError(http.StatusConflict, "Username taken"))
		} else {
			logger.WithError(err).Debug("Failed to register the user")
			ape.RenderErr(w, apierrors.NewApiError(
				http.StatusInternalServerError, "Failed to register the user",
			))
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
