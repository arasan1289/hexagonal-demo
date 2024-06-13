package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Logger struct {
	zerolog.Logger
	QueryThreshold uint
}

var (
	zerologLogger *Logger
)

// Set returns a logger instance
func Set(c *config.Container) (*Logger, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var level zerolog.Level
	var log zerolog.Logger
	switch c.App.Env {
	case "development", "dev":
		level = zerolog.DebugLevel
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("Message: %s", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		output.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		log = zerolog.New(output).With().Timestamp().Logger()
	case "production", "prod":
		level = zerolog.InfoLevel
		output := os.Stdout
		log = zerolog.New(output).With().Timestamp().Logger()
	default:
		return nil, fmt.Errorf("unknown environment: %s", c.App.Env)
	}
	zerolog.SetGlobalLevel(level)
	zerologLogger = &Logger{log, c.App.QueryThreshold}
	return zerologLogger, nil
}

// GetLogger returns the Zerolog logger.
func GetLogger() *Logger {
	return zerologLogger
}

type gormLogger struct {
	logger *Logger
}

// NewGormLogger returns a new GORM logger that logs to Zerolog.
func NewGormLogger() logger.Interface {
	return &gormLogger{logger: GetLogger()}
}

// LogMode implements the gorm.Logger interface.
func (gl *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := gl.logger
	switch level {
	case logger.Silent:
		newLogger.Level(zerolog.Disabled)
	case logger.Error:
		newLogger.Level(zerolog.ErrorLevel)
	case logger.Warn:
		newLogger.Level(zerolog.WarnLevel)
	case logger.Info:
		newLogger.Level(zerolog.InfoLevel)
	default:
		newLogger.Level(zerolog.DebugLevel)
	}
	return &gormLogger{
		logger: newLogger,
	}
}

// Trace implements the gorm.Logger interface.
func (gl *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// elapsed := time.Since(begin)

	// Log executed queries and record not found errors
	sql, rows := fc()
	duration := time.Since(begin).Milliseconds()

	// Set threshold for long-running queries (e.g., 100 milliseconds)
	threshold := int64(gl.logger.QueryThreshold) // in milliseconds

	if duration > threshold {
		// Log long-running query
		gl.logger.Warn().
			Str("module", "gorm").
			Str("sql", sql).
			Int64("rows", rows).
			Int64("duration_ms", duration).
			Msg("Long-running query")
	}
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			gl.logger.Error().
				Str("module", "gorm").
				Str("sql", sql).
				Int64("duration_ms", duration).
				Msg("Record not found")
		case gorm.ErrDuplicatedKey:
			gl.logger.Error().
				Str("module", "gorm").
				Int64("duration_ms", duration).
				Str("sql", sql).Err(err).
				Msg("Duplicate key")
		default:
			gl.logger.Error().
				Str("module", "gorm").
				Str("sql", sql).
				Int64("duration_ms", duration).
				Err(err).
				Msg("Query execution failed")
		}
	} else {
		gl.logger.Debug().
			Str("module", "gorm").
			Str("sql", sql).
			Int64("rows", rows).
			Int64("duration_ms", duration).
			Msg("Query executed")
	}
}

// Info implements the gorm.Logger interface.
func (gl *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	gl.logger.Info().Msgf(msg, data...)
}

// Warn implements the gorm.Logger interface.
func (gl *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	gl.logger.Warn().Msgf(msg, data...)
}

// Error implements the gorm.Logger interface.
func (gl *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	gl.logger.Error().Msgf(msg, data...)
}
