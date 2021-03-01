package log

import (
	"github.com/go-logr/logr"
	log "github.com/sirupsen/logrus"
)

type StaticLogger struct {
}

// Enabled tests whether this Logger is enabled.  For example, commandline
// flags might be used to set the logging verbosity and disable some info
// logs.
func (m *StaticLogger) Enabled() bool {
	return true
}

// Info logs a non-error message with the given key/value pairs as context.
//
// The msg argument should be used to add some constant description to
// the log line.  The key/value pairs can then be used to add additional
// variable information.  The key/value pairs should alternate string
// keys and arbitrary values.
func (m *StaticLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Infof(msg, keysAndValues...)
}

// Error logs an error, with the given message and key/value pairs as context.
// It functions similarly to calling Info with the "error" named value, but may
// have unique behavior, and should be preferred for logging errors (see the
// package documentations for more information).
//
// The msg field should be used to add context to any underlying error,
// while the err field should be used to attach the actual error that
// triggered this log line, if present.
func (m *StaticLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	log.WithFields(log.Fields{"error": err}).Errorf(msg, keysAndValues...)
}

// V returns an Logger value for a specific verbosity level, relative to
// this Logger.  In other words, V values are additive.  V higher verbosity
// level means a log message is less important.  It's illegal to pass a log
// level less than zero.
func (m *StaticLogger) V(level int) logr.Logger {
	return m
}

// WithValues adds some key-value pairs of context to a logger.
// See Info for documentation on how key/value pairs work.
func (m *StaticLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	return m
}

// WithName adds a new element to the logger's name.
// Successive calls with WithName continue to append
// suffixes to the logger's name.  It's strongly reccomended
// that name segments contain only letters, digits, and hyphens
// (see the package documentation for more information).
func (m *StaticLogger) WithName(name string) logr.Logger {
	return NewLogger(log.WithField("name", name))
}

type Logger struct {
	logger  log.FieldLogger
	enabled bool
}

func NewLogger(logger log.FieldLogger) *Logger {
	return &Logger{
		logger:  logger,
		enabled: true,
	}
}

// Enabled tests whether this Logger is enabled.  For example, commandline
// flags might be used to set the logging verbosity and disable some info
// logs.
func (l *Logger) Enabled() bool {
	return l.enabled
}

// Info logs a non-error message with the given key/value pairs as context.
//
// The msg argument should be used to add some constant description to
// the log line.  The key/value pairs can then be used to add additional
// variable information.  The key/value pairs should alternate string
// keys and arbitrary values.
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Infof(msg, keysAndValues...)
}

// Error logs an error, with the given message and key/value pairs as context.
// It functions similarly to calling Info with the "error" named value, but may
// have unique behavior, and should be preferred for logging errors (see the
// package documentations for more information).
//
// The msg field should be used to add context to any underlying error,
// while the err field should be used to attach the actual error that
// triggered this log line, if present.
func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.WithFields(log.Fields{"error": err}).Errorf(msg, keysAndValues...)
}

// V returns an Logger value for a specific verbosity level, relative to
// this Logger.  In other words, V values are additive.  V higher verbosity
// level means a log message is less important.  It's illegal to pass a log
// level less than zero.
func (l *Logger) V(level int) logr.Logger {
	return l
}

// WithValues adds some key-value pairs of context to a logger.
// See Info for documentation on how key/value pairs work.
func (l *Logger) WithValues(keysAndValues ...interface{}) logr.Logger {
	return l
}

// WithName adds a new element to the logger's name.
// Successive calls with WithName continue to append
// suffixes to the logger's name.  It's strongly reccomended
// that name segments contain only letters, digits, and hyphens
// (see the package documentation for more information).
func (l *Logger) WithName(name string) logr.Logger {
	return NewLogger(log.WithField("name", name))
}
