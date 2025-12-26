package handlers

import (
	"fmt"
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/service/errors/apierrors"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
)

func Login(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)
	request, err := requests.NewLoginRequest(r)
	if err != nil {
		logger.WithError(err).Debug("bad request")
		ape.RenderErr(w, apierrors.BadRequest())
	}
	ape.Render(w, fmt.Sprintf("WIP; Request: %v", request))
}
