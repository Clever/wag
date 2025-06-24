package wagclientlogger

type NoOpLogger struct{}

// NoOpLogger.Log does not have a function body, since it's intended as a mock logger solely to
// satisfy the interface in tests where we are not asserting on logging behavior.
func (l NoOpLogger) Log(level LogLevel, message string, pairs map[string]interface{}) {}
