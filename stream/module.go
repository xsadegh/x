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
	Creds   string `yaml:"creds"`
}

func New(config Config, logger *zap.Logger) jetstream.JetStream {
	var err error
	var nc *nats.Conn
	var js jetstream.JetStream

	var opts = []nats.Option{nats.Name(config.Client)}

	if config.Token != "" {
		opts = append(opts, nats.Token(config.Token))
	}

	if config.Creds != "" {
		opts = append(opts, nats.UserCredentials(config.Creds))
	}

	nc, err = nats.Connect(config.Address, opts...)
	if err != nil {
		logger.Fatal("failed to connect to nats", zap.Error(err))
	}

	js, err = jetstream.New(nc)
	if err != nil {
		logger.Fatal("failed to initialize jetstream", zap.Error(err))
	}

	return js
}
