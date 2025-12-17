package store

import (
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var MODULE = fx.Module(
	"STORE",
	fx.Provide(New),
)

type Config struct {
	Endpoints []string `yaml:"endpoints"`
}

func New(config Config, logger *zap.Logger) *etcd.Client {
	store, err := etcd.New(etcd.Config{Endpoints: config.Endpoints})
	if err != nil {
		logger.Fatal("failed to connect to store", zap.Error(err))
	}

	return store
}
