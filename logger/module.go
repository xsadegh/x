package logger

import (
	"os"
	"strings"

	"go.elastic.co/apm/module/apmzap/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var MODULE = fx.Module(
	"LOGGER",
	fx.Provide(NewLogger),
)

type Config struct {
	APM      bool   `yaml:"-"`
	Debug    bool   `yaml:"debug"`
	Encode   string `yaml:"encode"`
	LogLevel string `yaml:"loglevel"`
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
	if err := level.Set(config.LogLevel); err != nil {
		_ = level.Set(zapcore.DebugLevel.String())
	}

	options := []zap.Option{zap.AddCaller()}
	if config.Debug {
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	}
	var core = zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zap.NewAtomicLevelAt(level))
	if config.APM {
		apm := apmzap.Core{}
		core = apm.WrapCore(core)
	}
	return zap.New(core, options...)
}
