package logger

import (
	"os"
	"strings"

	adapter "github.com/axiomhq/axiom-go/adapters/zap"
	"github.com/axiomhq/axiom-go/axiom"
	"github.com/axiomhq/axiom-go/axiom/ingest"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.sadegh.io/x/tracer"
)

var MODULE = fx.Module(
	"LOGGER",
	fx.Provide(NewLogger),
)

type Config struct {
	Debug  bool   `yaml:"debug"`
	Level  string `yaml:"level"`
	Encode string `yaml:"encode"`

	Tracer tracer.Config `yaml:"-"`
}

func Logger(config Config, logger *zap.Logger) fxevent.Logger {
	if config.Debug {
		return &fxevent.ZapLogger{Logger: logger}
	}

	return fxevent.NopLogger
}

func NewLogger(config Config) *zap.Logger {
	var encoderConfig zapcore.EncoderConfig
	if config.Debug {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	var encoder zapcore.Encoder
	switch strings.ToLower(config.Encode) {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var level zapcore.Level
	if err := level.Set(config.Level); err != nil {
		_ = level.Set(zapcore.DebugLevel.String())
	}

	token := strings.ReplaceAll(config.Tracer.Headers["Authorization"], "Bearer ", "")
	options := []zap.Option{zap.AddCaller()}
	if config.Debug {
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	}
	options = append(options)

	var core = zapcore.NewCore(
		encoder,
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(level),
	)

	var dupCore, _ = adapter.New(
		adapter.SetDataset(config.Tracer.Logs),
		adapter.SetClientOptions(axiom.SetToken(token)),
		adapter.SetIngestOptions(
			ingest.SetEventLabel("service", config.Tracer.Service),
			ingest.SetEventLabel("version", config.Tracer.Version),
		),
	)

	if config.Tracer.Logs == "" {
		return zap.New(core, options...)
	} else {
		return zap.New(zapcore.NewTee(core, dupCore), options...)
	}
}
