package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug   LogLevel = "DEBUG"
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
)

// LoggerConfig holds logging configuration
type LoggerConfig struct {
	Level      LogLevel
	LogFile    string
	MaxSize    int  // megabytes
	MaxBackups int  // number of old log files to retain
	MaxAge     int  // days
	Compress   bool // compress rotated files
}

// DefaultLoggerConfig returns default logging configuration
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:      LogLevelInfo,
		LogFile:    "./logs/bot.log",
		MaxSize:    10,   // 10 MB
		MaxBackups: 5,    // keep 5 old log files
		MaxAge:     30,   // keep logs for 30 days
		Compress:   true, // compress old logs
	}
}

// SetupLogger initializes the logging infrastructure with file rotation
func SetupLogger(config LoggerConfig) *slog.Logger {
	// Parse log level
	level := parseLogLevel(config.Level)

	// Create logs directory if it doesn't exist
	if config.LogFile != "" {
		logDir := filepath.Dir(config.LogFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			slog.Error("Failed to create log directory", "error", err, "dir", logDir)
		}
	}

	// Create multi-writer: stdout + file with rotation
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if config.LogFile != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		writers = append(writers, fileWriter)
	}

	multiWriter := io.MultiWriter(writers...)

	// Create JSON handler with appropriate log level
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: level,
		AddSource: level == slog.LevelDebug, // add source location for debug logs
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

// parseLogLevel converts string log level to slog.Level
func parseLogLevel(level LogLevel) slog.Level {
	switch level {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarning:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// GetLoggerFromEnv creates logger configuration from environment variables
func GetLoggerFromEnv() LoggerConfig {
	config := DefaultLoggerConfig()

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Level = LogLevel(level)
	}

	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		config.LogFile = logFile
	}

	return config
}
