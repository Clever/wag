package wagclientlogger

// WagClientLogger provides a minimal interface for a Wag Client Logger
type WagClientLogger interface {
	Log(level LogLevel, message string, pairs map[string]interface{})
}

// M is a convenience type to avoid having to type out map[string]interface{} everytime.
type M map[string]interface{}

// LogLevel is an enum of valid log levels.
type LogLevel int

// Constants used to define different LogLevels supported
const (
	Trace LogLevel = iota
	Debug
	Info
	Warning
	Error
	Critical
	FromEnv
)
