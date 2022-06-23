package logger

//WagClientLogger provides a minimal interface for a Wag Logger
type WagClientLogger interface {
	Log(level int, message string, pairs map[string]interface{})
}
