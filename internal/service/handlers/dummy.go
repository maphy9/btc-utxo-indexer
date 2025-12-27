package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
)

func Dummy(w http.ResponseWriter, r *http.Request) {
	userID := helpers.UserID(r)
	db := helpers.DB(r)

	user, err := db.Users().GetByUserID(r.Context(), userID)
	if err != nil {
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError, "Could not find the user",
		))
		return
	}

	payload, err := json.Marshal(user)
	if err != nil {
		ape.RenderErr(w, apierrors.NewApiError(
			http.StatusInternalServerError, "Marshal failed",
		))
		return
	}
	ape.Render(w, string(payload))
}
