package service

import (
	"github.com/go-chi/chi"
	"github.com/maphy9/btc-utxo-indexer/internal/service/handlers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/helpers"
	"github.com/maphy9/btc-utxo-indexer/internal/service/middleware"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
			helpers.CtxServiceConfig(s.serviceConfig),
			helpers.CtxDB(s.db),
			helpers.CtxManager(s.manager),
		),
	)

	r.Post("/login", handlers.Login)
	r.Post("/register", handlers.Register)
	r.Post("/refresh", handlers.Refresh)

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.AuthMiddleware,
		)

		r.Route("/addresses", func(r chi.Router) {
			r.Post("/", handlers.AddAddress)
			r.Get("/", handlers.GetAddresses)
			r.Get("/{address}/utxos", handlers.GetUtxos)
		})
	})

	return r
}
