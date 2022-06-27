package printlogger

import (
	"fmt"

	"github.com/Clever/wag/logger"
)

//NewLogger creates a logger for id that produces logs at and below the indicated level.
func NewLogger(id string, level string) logger.WagClientLogger {
	return PrintlnLogger{id: id, level: level}
}

type PrintlnLogger struct {
	level string
	id    string
}

func (w PrintlnLogger) Log(level string, message string, m map[string]interface{}) {

	if w.strLvlToInt(w.level) >= w.strLvlToInt(level) {
		fmt.Print(w.id, ": ")
		fmt.Print(message)
		for k, v := range m {
			fmt.Print(" ", k, " : ", v)
		}
		fmt.Println()
	}
}

func (w PrintlnLogger) strLvlToInt(s string) int {
	switch s {
	case "Critical":
		return 0
	case "Error":
		return 1
	case "Warning":
		return 2
	case "Info":
		return 3
	case "Debug":
		return 4
	case "Trace":
		return 5
	}
	return -1
}
