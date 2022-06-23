package printlnLoggerInterface

import (
	"fmt"

	"github.com/Clever/wag/logger"
)

//NewLogger creates a logger for id that produces logs at and below the indicated level.
func NewLogger(id string, level int) logger.WagClientLogger {
	return PrintlnLogger{id: id, level: level}
}

type PrintlnLogger struct {
	level int
	id    string
}

func (w PrintlnLogger) Log(level int, message string, m map[string]interface{}) {
	if w.level >= level {
		fmt.Print(w.id, ": ")
		fmt.Print(message)
		for k, v := range m {
			fmt.Print(" ", k, " : ", v)
		}
		fmt.Println()
	}
}
