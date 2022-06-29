package printlogger

import (
	"encoding/json"
	"fmt"

	logger "github.com/Clever/wag/loggers/waglogger"
)

//NewLogger creates a logger for id that produces logs at and below the indicated level.
//Level indicated the level at and below which logs are created.
func NewLogger(id string, level string) logger.WagClientLogger {
	return PrintlnLogger{id: id, level: level}
}

type PrintlnLogger struct {
	level string
	id    string
}

func (w PrintlnLogger) Log(level string, message string, m map[string]interface{}) {

	if w.strLvlToInt(w.level) >= w.strLvlToInt(level) {
		m["id"] = w.id
		jsonLog, err := json.Marshal(m)
		if err != nil {
			jsonLog, err = json.Marshal(map[string]interface{}{"Error Marshalling Log": err})
		}
		fmt.Println(string(jsonLog))
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
