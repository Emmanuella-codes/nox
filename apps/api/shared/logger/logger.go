package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(cfg *config.Config, service string) error {
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	writer := stdoutWriter(cfg)
	if shouldWriteLocalFile(cfg) {
		fileWriter, err := NewRotatingFileWriter(RotatingFileConfig{
			Path:       cfg.LogFilePath,
			MaxSizeMB:  cfg.LogFileMaxSizeMB,
			MaxBackups: cfg.LogFileMaxBackups,
			MaxAgeDays: cfg.LogFileMaxAgeDays,
			Compress:   cfg.LogFileCompress,
		})
		if err != nil {
			return err
		}
		writer = io.MultiWriter(writer, fileWriter)
	}

	log.Logger = zerolog.New(writer).
		With().
		Timestamp().
		Str("service", service).
		Str("env", cfg.Environment).
		Logger()
	return nil
}

func stdoutWriter(cfg *config.Config) io.Writer {
	if strings.EqualFold(cfg.LogFormat, "console") {
		return zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}
	return os.Stdout
}

func shouldWriteLocalFile(cfg *config.Config) bool {
	return cfg.LogFileEnabled && strings.EqualFold(cfg.Environment, "development")
}
