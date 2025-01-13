package log

import (
	"app/config"
	"log/slog"
	"os"

	"go.elastic.co/ecszap"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

// LogService struct
type LogService struct {
	*slog.Logger
}

func newService(conf *config.Config) *LogService {
	var (
		log  *slog.Logger
		zLog *zap.Logger
	)
	if conf.AppEnv == "production" {
		zapOption := []zap.Option{
			zap.AddCaller(),
		}
		if conf.Debug {
			zapOption = append(zapOption, zap.Development())
		}

		encoderConfig := ecszap.NewDefaultEncoderConfig()
		core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
		zLog = zap.New(core, zapOption...)
	} else {
		var err error
		zLog, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}
	core := zapcore.NewTee(
		zLog.Core(),
		otelzap.NewCore(conf.AppName),
	)
	log = slog.New(zapslog.NewHandler(core, &zapslog.HandlerOptions{
		AddSource:  true,
		LoggerName: conf.AppName,
	}))
	slog.SetDefault(log)

	return &LogService{
		Logger: log,
	}
}
