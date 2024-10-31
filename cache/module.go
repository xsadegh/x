package cache

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var MODULE = fx.Module(
	"CACHE",
	fx.Provide(New),
)

type Config struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	TLSMode  bool   `yaml:"tlsMode"`
}

func New(config Config) *redis.Client {
	var opts = &redis.Options{
		Addr: config.Address, Password: config.Password,
		ConnMaxIdleTime: -1, MaxRetries: 10, PoolSize: 100,
	}

	opts.PoolTimeout = 2 * time.Minute
	opts.ReadTimeout = 2 * time.Minute
	opts.WriteTimeout = 1 * time.Minute

	if config.TLSMode {
		host, _, _ := net.SplitHostPort(config.Address)
		opts.TLSConfig = &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}
	}

	return redis.NewClient(opts)
}
