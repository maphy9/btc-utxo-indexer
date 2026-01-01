package service

import (
	"net"
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain"
	"github.com/maphy9/btc-utxo-indexer/internal/config"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log           *logan.Entry
	copus         types.Copus
	listener      net.Listener
	serviceConfig *config.ServiceConfig
	db            *pgdb.DB
	manager       *blockchain.Manager
}

func (s *service) run() error {
	s.log.Info("Service started")
	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) (*service, error) {
	manager, err := blockchain.NewManager("electrum.blockstream.info:50002")
	if err != nil {
		return nil, err
	}
	return &service{
		log:           cfg.Log(),
		copus:         cfg.Copus(),
		listener:      cfg.Listener(),
		serviceConfig: cfg.ServiceConfig(),
		db:            cfg.DB(),
		manager:       manager,
	}, nil
}

func Run(cfg config.Config) {
	service, err := newService(cfg)
	if err != nil {
		panic(err)
	}
	if err := service.run(); err != nil {
		panic(err)
	}
}
