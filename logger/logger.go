package logger

// Logger defines an interface to a logger.
type Logger interface {
	Debug() Event
	Info() Event
	Warn() Event
	Error() Event
	Fatal() Event
	Panic() Event
}

// Event defines an interface to a log event containing additional context
// for a log message. The Msg (or Send) method must be called on the event
// in order to send it.
type Event interface {
	Send()
	Msg(string)
	Str(string, string) Event
	Int(string, int) Event
	Uint(string, uint) Event
	Err(error) Event
}
