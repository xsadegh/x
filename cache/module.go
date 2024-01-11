package cache

import (
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
}

func New(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		ConnMaxIdleTime: -1, MaxRetries: 10, PoolSize: 100,
		Addr: config.Address, Password: config.Password,
		PoolTimeout:  2 * time.Minute,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	})
}
