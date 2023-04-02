package groupcache

import (
	"github.com/mailgun/groupcache/v2"

	"darvaza.org/slog"
)

var (
	_ groupcache.Logger = (*Logger)(nil)
)

// Logger is a specific log context for groupcache
type Logger struct {
	logger slog.Logger
}

// Printf logs a message under a previously set level and with previously set fields
func (gcl *Logger) Printf(format string, args ...any) {
	gcl.logger.Printf(format, args...)
}

// Error creates a new logger context with level set to Error
func (gcl *Logger) Error() groupcache.Logger {
	return &Logger{
		logger: gcl.logger.Error(),
	}
}

// Warn creates a new logger context with level set to Warning
func (gcl *Logger) Warn() groupcache.Logger {
	return &Logger{
		logger: gcl.logger.Warn(),
	}
}

// Info creates a new logger context with level set to Info
func (gcl *Logger) Info() groupcache.Logger {
	return &Logger{
		logger: gcl.logger.Info(),
	}
}

// Debug creates a new logger context with level set to Debug
func (gcl *Logger) Debug() groupcache.Logger {
	return &Logger{
		logger: gcl.logger.Debug(),
	}
}

// ErrorField creates a new logger context with a new field containing an error
func (gcl *Logger) ErrorField(label string, err error) groupcache.Logger {
	return &Logger{
		logger: gcl.logger.WithField(label, err),
	}
}

// StringField creates a new logger context with a new field containing a string value
func (gcl *Logger) StringField(label string, val string) groupcache.Logger {
	return &Logger{
		logger: gcl.logger.WithField(label, val),
	}
}

// WithFields creates a new logger context with a set of new fields of arbitrary value
func (gcl *Logger) WithFields(fields map[string]any) groupcache.Logger {
	return &Logger{
		logger: gcl.logger.WithFields(fields),
	}
}

// NewLogger creates a Logger for groupcache wrapping a given slog.Logger
func NewLogger(l slog.Logger) groupcache.Logger {
	return &Logger{
		logger: l,
	}
}

// SetLogger sets groupcache to use a given slog.Logger
func SetLogger(l slog.Logger) {
	gcl := NewLogger(l)
	groupcache.SetLoggerFromLogger(gcl)
}
