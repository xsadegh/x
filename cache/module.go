package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var MODULE = fx.Module(
	"CACHE",
	fx.Provide(NewPool),
)

type Config struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

func NewPool(config *Config, logger *zap.Logger) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			var opts []redis.DialOption
			if config.Password != "" {
				opts = append(opts, redis.DialPassword(config.Password))
			}

			conn, err := redis.Dial("tcp", config.Address, opts...)
			if err != nil {
				logger.Fatal("failed to dail redis", zap.Error(err))
			}

			return conn, err
		},
		MaxIdle:   50,
		MaxActive: 10000,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
