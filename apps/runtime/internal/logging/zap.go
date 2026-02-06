package logging

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.SugaredLogger for structured logging with backward compatibility
type Logger struct {
	*zap.SugaredLogger
	mu sync.RWMutex
}

// Global logger instance
var (
	defaultLogger *Logger
	initOnce      sync.Once
)

// Config holds logger configuration
type Config struct {
	Level      string // debug, info, warn, error
	Output     string // stdout, stderr, or file path
	Format     string // json, console
	TimeFormat string
	Prefix     string
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		Output:     "stdout",
		Format:     "console",
		TimeFormat: time.RFC3339,
		Prefix:     "",
	}
}

// NewLogger creates a new Logger with the given configuration
func NewLogger(cfg Config) (*Logger, error) {
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn", "warning":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	var encoding string
	if cfg.Format == "json" {
		encoding = "json"
	} else {
		encoding = "console"
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Encoding:    encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{cfg.Output},
		ErrorOutputPaths: []string{cfg.Output},
		InitialFields:    map[string]interface{}{},
	}

	if cfg.Prefix != "" {
		config.InitialFields["prefix"] = cfg.Prefix
	}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}, nil
}

// Default returns the default global logger, initializing it if necessary
func Default() *Logger {
	initOnce.Do(func() {
		var err error
		defaultLogger, err = NewLogger(DefaultConfig())
		if err != nil {
			// Fallback to standard logger if zap fails
			defaultLogger = &Logger{
				SugaredLogger: zap.NewNop().Sugar(),
			}
		}
	})
	return defaultLogger
}

// SetDefault sets the default global logger
func SetDefault(logger *Logger) {
	defaultLogger = logger
}

// Sync flushes any buffered log entries
func Sync() {
	if defaultLogger != nil && defaultLogger.SugaredLogger != nil {
		defaultLogger.Sync()
	}
}

// With creates a child logger with additional context
func With(key string, value interface{}) *Logger {
	Default().mu.Lock()
	defer Default().mu.Unlock()
	return &Logger{
		SugaredLogger: Default().With(key, value),
	}
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	Default().Debug(args...)
}

// Debugf logs a formatted debug message
func Debugf(template string, args ...interface{}) {
	Default().Debugf(template, args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	Default().Info(args...)
}

// Infof logs a formatted info message
func Infof(template string, args ...interface{}) {
	Default().Infof(template, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	Default().Warn(args...)
}

// Warnf logs a formatted warning message
func Warnf(template string, args ...interface{}) {
	Default().Warnf(template, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	Default().Error(args...)
}

// Errorf logs a formatted error message
func Errorf(template string, args ...interface{}) {
	Default().Errorf(template, args...)
}

// Fatal logs a fatal message and calls os.Exit(1)
func Fatal(args ...interface{}) {
	Default().Fatal(args...)
}

// Fatalf logs a formatted fatal message and calls os.Exit(1)
func Fatalf(template string, args ...interface{}) {
	Default().Fatalf(template, args...)
}

// Print logs a message (compatible with standard log)
func Print(args ...interface{}) {
	Default().Info(args...)
}

// Printf logs a formatted message (compatible with standard log)
func Printf(template string, args ...interface{}) {
	Default().Infof(template, args...)
}

// Println logs a message with newline (compatible with standard log)
func Println(args ...interface{}) {
	Default().Info(args...)
}

// Panic logs a panic message and panics (compatible with standard log)
func Panic(args ...interface{}) {
	Default().Panic(args...)
}

// Panicln logs a panic message with newline and panics (compatible with standard log)
func Panicln(args ...interface{}) {
	Default().Panic(args...)
}
