package service

import (
	"github.com/go-chi/chi"
	"github.com/maphy9/btc-utxo-indexer/internal/service/handlers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
		),
	)
	r.Route("/", func(r chi.Router) {
		r.Post("/login", handlers.Login)
	})

	return r
}
