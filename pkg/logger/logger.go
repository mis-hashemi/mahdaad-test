package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	FilePath         string
	UseLocalTime     bool
	FileMaxSizeInMB  int
	FileMaxAgeInDays int
	Level            slog.Level
}

var defaultConfig = Config{
	UseLocalTime:     false,
	FileMaxSizeInMB:  10,
	FileMaxAgeInDays: 30,
	Level:            slog.LevelInfo,
}

type Logger struct {
	*slog.Logger
}

func New(cfg *Config) *Logger {
	if cfg == nil {
		cfg = &defaultConfig
	}

	writer := &lumberjack.Logger{
		Filename:  cfg.FilePath,
		LocalTime: cfg.UseLocalTime,
		MaxSize:   cfg.FileMaxSizeInMB,
		MaxAge:    cfg.FileMaxAgeInDays,
	}

	opts := &slog.HandlerOptions{Level: cfg.Level}

	base := slog.New(
		slog.NewJSONHandler(io.MultiWriter(writer, os.Stdout), opts),
	)

	return &Logger{Logger: base}
}

func (l *Logger) WithCaller() *Logger {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}
	return &Logger{Logger: l.With("caller", fmt.Sprintf("%s:%d", file, line))}
}


