package logger

//WagClientLogger provides a minimal interface for a Wag Logger
type WagClientLogger interface {
	Log(level string, message string, pairs map[string]interface{})
}

//M is a convenience type to avoid having to type out map[string]interface{} everytime.
type M map[string]interface{}
