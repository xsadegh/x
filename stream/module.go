package stream

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var MODULE = fx.Module(
	"STREAM",
	fx.Provide(New),
)

type Config struct {
	Address string `yaml:"address"`
	Client  string `yaml:"client"`
	Token   string `yaml:"token"`
}

func New(config Config, logger *zap.Logger) jetstream.JetStream {
	var err error
	var nc *nats.Conn
	var js jetstream.JetStream
	nc, err = nats.Connect(
		config.Address,
		nats.MaxReconnects(-1),
		nats.Name(config.Client),
		nats.Token(config.Token),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Debug("stream disconnected", zap.Error(err))
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Debug("stream reconnected")
		}),
	)
	if err != nil {
		logger.Fatal("failed to connect to nats", zap.Error(err))
	}

	js, err = jetstream.New(nc)
	if err != nil {
		logger.Fatal("failed to initialize jetstream", zap.Error(err))
	}

	return js
}
