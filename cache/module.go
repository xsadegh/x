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
	TLSMode  bool   `yaml:"tlsMode"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func New(config Config) *redis.Client {
	var opts = &redis.Options{}

	opts.MaxRetries = 0
	opts.Addr = config.Address
	opts.Username = config.Username
	opts.Password = config.Password
	opts.PoolTimeout = 2 * time.Minute
	opts.ReadTimeout = 2 * time.Minute
	opts.WriteTimeout = 1 * time.Minute

	if config.TLSMode {
		host, _, _ := net.SplitHostPort(config.Address)
		opts.TLSConfig = &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}
	}

	return redis.NewClient(opts)
}
