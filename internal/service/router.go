package service

import (
  "https://github.com/maphy9/btc-utxo-indexer/internal/service/handlers"
  "github.com/go-chi/chi"
  "gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
  r := chi.NewRouter()

  r.Use(
    ape.RecoverMiddleware(s.log),
    ape.LoganMiddleware(s.log),
    ape.CtxMiddleware(
      handlers.CtxLog(s.log),
    ),
  )
  r.Route("/integrations/btc-utxo-indexer", func(r chi.Router) {
    // configure endpoints here
  })

  return r
}
