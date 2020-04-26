package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// zerologLogger is a logger implementing the Logger implementation using zerolog as the underlying logging mechanism.
type zerologLogger struct {
	logger *zerolog.Logger
}

// zerologEvent wraps a zerolog.Event and implements the Event interface.
type zerologEvent struct {
	event *zerolog.Event
}

// NewZerologLogger initializes and returns a new ZerologLogger.
func NewZerologLogger(isDebug bool) Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &zerologLogger{logger: &logger}
}

// Debug starts a new log message with the log level set to debug.
// The Msg method must be called on the returned event in order to send it.
func (l *zerologLogger) Debug() Event {
	return &zerologEvent{event: l.logger.Debug()}
}

// Info starts a new log message with the log level set to info.
// The Msg method must be called on the returned event in order to send it.
func (l *zerologLogger) Info() Event {
	return &zerologEvent{event: l.logger.Info()}
}

// Warn starts a new log message with the log level set to warn.
// The Msg method must be called on the returned event in order to send it.
func (l *zerologLogger) Warn() Event {
	return &zerologEvent{event: l.logger.Warn()}
}

// Error starts a new log message with the log level set to error.
// The Msg method must be called on the returned event in order to send it.
func (l *zerologLogger) Error() Event {
	return &zerologEvent{event: l.logger.Error()}
}

// Fatal starts a new log message with the log level set to fatal.
// The os.Exit(1) function is called by the Msg method, which terminates the program immediately.
// The Msg method must be called on the returned event in order to send it.
func (l *zerologLogger) Fatal() Event {
	return &zerologEvent{event: l.logger.Fatal()}
}

// Panic starts a new log message with the log level set to panic.
// The panic() function is called by the Msg method.
// The Msg method must be called on the returned event in order to send it.
func (l *zerologLogger) Panic() Event {
	return &zerologEvent{event: l.logger.Panic()}
}

// Msg sends the event with msg added as the message field if not empty.
func (e *zerologEvent) Msg(msg string) {
	e.event.Msg(msg)
}

// Send is equivalent to calling Msg("").
func (e *zerologEvent) Send() {
	e.event.Send()
}

// Str adds the field key with str as a string to the Event context.
func (e *zerologEvent) Str(key string, str string) Event {
	return &zerologEvent{event: e.event.Str(key, str)}
}

// Int adds the field key with num as a int to the Event context.
func (e *zerologEvent) Int(key string, num int) Event {
	return &zerologEvent{event: e.event.Int(key, num)}
}

// Uint adds the field key with num as a uint to the Event context.
func (e *zerologEvent) Uint(key string, num uint) Event {
	return &zerologEvent{event: e.event.Uint(key, num)}
}

// Err adds the field "error" with serialized err to the Event context.
// If err is nil, no field is added.
func (e *zerologEvent) Err(err error) Event {
	return &zerologEvent{event: e.event.Err(err)}
}
