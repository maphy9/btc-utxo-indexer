package config

import (
	"time"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

func NewServiceConfiger(getter kv.Getter) ServiceConfiger {
	return &serviceConfiger{
		getter: getter,
	}
}

type ServiceConfiger interface {
	ServiceConfig() *ServiceConfig
}

type ServiceConfig struct {
	AccessTokenKey         string        `fig:"access_token_key,required"`
	AccessTokenExpireTime  time.Duration `fig:"access_token_expire_time,required"`
	RefreshTokenKey        string        `fig:"refresh_token_key,required"`
	RefreshTokenExpireTime time.Duration `fig:"refresh_token_expire_time,required"`
	RawNodes               map[string]struct {
		SSL               bool   `fig:"ssl"`
		ReconnectAttempts uint32 `fig:"reconnect_attempts"`
	} `fig:"nodes"`
	NodeEntries []blockchain.NodepoolEntry
}

type serviceConfiger struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *serviceConfiger) ServiceConfig() *ServiceConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "service")
		config := ServiceConfig{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic("Failed to read service config")
		}
		config.NodeEntries = make([]blockchain.NodepoolEntry, 0, len(config.RawNodes))
		for address, nodeCfg := range config.RawNodes {
			entry := blockchain.NodepoolEntry{
				Address:           address,
				SSL:               nodeCfg.SSL,
				ReconnectAttempts: nodeCfg.ReconnectAttempts,
			}
			config.NodeEntries = append(config.NodeEntries, entry)
		}
		return &config
	}).(*ServiceConfig)
}
