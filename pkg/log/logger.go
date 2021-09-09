package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Logger      *zap.SugaredLogger
	debugOption bool
)

func InitLogger(debug, quiet bool) (err error) {
	debugOption = debug
	Logger, err = newLogger(debug, quiet)
	if err != nil {
		return fmt.Errorf("error in new logger: %w", err)
	}
	return nil
}

func newLogger(debug, quiet bool) (*zap.SugaredLogger, error) {
	level := zap.NewAtomicLevel()
	if debug {
		level.SetLevel(zapcore.DebugLevel)
	} else {
		level.SetLevel(zapcore.InfoLevel)
	}

	stdout := "stdout"
	stderr := "stderr"
	if quiet {
		if _, err := os.Create(os.DevNull); err != nil {
			return nil, err
		}
		stdout = os.DevNull
		stderr = os.DevNull
	}

	myConfig := zap.Config{
		Level:             level,
		Encoding:          "console",
		Development:       debug,
		DisableStacktrace: true,
		DisableCaller:     true,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{stdout},
		ErrorOutputPaths: []string{stderr},
	}
	logger, err := myConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap config: %w", err)
	}

	return logger.Sugar(), nil
}

func Fatal(err error) {
	if debugOption {
		Logger.Fatalf("%+v", err)
	}
	Logger.Fatal(err)
}
