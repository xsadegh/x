package cache

import (
	"crypto/tls"
	"net"

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

	opts.Addr, opts.MaxRetries = config.Address, -1
	opts.Username, opts.Password = config.Username, config.Password

	if config.TLSMode {
		host, _, _ := net.SplitHostPort(config.Address)
		opts.TLSConfig = &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}
	}

	return redis.NewClient(opts)
}
