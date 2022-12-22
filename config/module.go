package config

import (
	"log"
	"os"

	"go.uber.org/config"
)

type Config interface {
	Path() string
	Env() map[string]any
}

func NewConfig[T any](base Config) T {
	var sources = []config.YAMLOption{
		config.Static(base.Env()), config.Expand(os.LookupEnv),
	}
	if _, err := os.Open(base.Path()); !os.IsNotExist(err) {
		sources = append(sources, config.File(base.Path()))
	}

	provider, err := config.NewYAML(sources...)
	if err != nil {
		log.Fatal(err)
	}

	var cfg T
	err = provider.Get(config.Root).Populate(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
