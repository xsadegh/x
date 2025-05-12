package stream

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

var MODULE = fx.Module(
	"STREAM",
	fx.Provide(New),
)

type Config struct {
	Credentials string `yaml:"credentials"`
	Address     string `yaml:"address"`
	Client      string `yaml:"client"`
}

func New(config Config) jetstream.JetStream {
	nc, _ := nats.Connect(
		config.Address,
		nats.Name(config.Client),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second),
		nats.UserCredentials(config.Credentials),
	)

	js, _ := jetstream.New(nc)

	return js
}
