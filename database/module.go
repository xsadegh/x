package database

import (
	"github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
)

var MODULE = fx.Module(
	"DATABASE",
	fx.Provide(NewDSN),
)

type Config struct {
	Debug    bool   `yaml:"debug"`
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type DSN string

func NewDSN(config Config) DSN {
	dsn := mysql.NewConfig()
	dsn.Net = "tcp"
	dsn.Addr = config.Address
	dsn.User = config.Username
	dsn.Passwd = config.Password
	dsn.DBName = config.Database
	dsn.Params = map[string]string{"parseTime": "true"}

	return DSN(dsn.FormatDSN())
}
