package helpers

import (
	"context"
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/config"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	serviceConfigCtxKey
	dbCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func CtxServiceConfig(serviceConfig *config.ServiceConfig) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, serviceConfigCtxKey, serviceConfig)
	}
}

func CtxDB(master data.MasterQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, dbCtxKey, master)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func ServiceConfig(r *http.Request) *config.ServiceConfig {
	return r.Context().Value(serviceConfigCtxKey).(*config.ServiceConfig)
}

func DB(r *http.Request) data.MasterQ {
	return r.Context().Value(dbCtxKey).(data.MasterQ)
}