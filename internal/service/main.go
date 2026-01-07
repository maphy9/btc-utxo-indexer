package service

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain"
	"github.com/maphy9/btc-utxo-indexer/internal/config"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"github.com/maphy9/btc-utxo-indexer/internal/data/pg"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log           *logan.Entry
	copus         types.Copus
	listener      net.Listener
	serviceConfig *config.ServiceConfig
	db            data.MasterQ
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
	db := pg.NewMasterQ(cfg.DB())
	log := cfg.Log()

	manager, err := blockchain.NewManager(cfg.ServiceConfig().NodeEntries, db, log)
	if err != nil {
		return nil, err
	}
	go manager.ListenHeaders()
	err = manager.SubscribeSavedAddresses()
	if err != nil {
		return nil, err
	}

	return &service{
		log:           log,
		copus:         cfg.Copus(),
		listener:      cfg.Listener(),
		serviceConfig: cfg.ServiceConfig(),
		db:            db,
		manager:       manager,
	}, nil
}

func Run(cfg config.Config) {
	service, err := newService(cfg)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := service.run(); err != nil {
			panic(err)
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	err = service.manager.Close()
	if err != nil {
		panic(err)
	}
}
