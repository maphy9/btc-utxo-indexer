package config

import (
	"time"

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
	TokenKey               string        `fig:"token_key,required"`
	TokenExpireTime        time.Duration `fig:"token_expire_time,required"`
	RefreshTokenKey        string        `fig:"refresh_token_key,required"`
	RefreshTokenExpireTime time.Duration `fig:"refresh_token_expire_time,required"`
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
		return &config
	}).(*ServiceConfig)
}
