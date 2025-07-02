package engines

import (
	"encoding/json"
	"os"
	"runtime/debug"
	"time"
)

type LoggerEngineFile struct {
	Path string
}

func (l *LoggerEngineFile) Trace(args ...interface{}) {
	filepath := l.Path + "/trace.log"
	writeLog(filepath, args, "Trace")
}

func (l *LoggerEngineFile) Debug(args ...interface{}) {
	filepath := l.Path + "/debug.log"
	writeLog(filepath, args, "Debug")
}

func (l *LoggerEngineFile) Info(args ...interface{}) {
	filepath := l.Path + "/info.log"
	writeLog(filepath, args, "Info")
}

func (l *LoggerEngineFile) Warn(args ...interface{}) {
	filepath := l.Path + "/warn.log"
	writeLog(filepath, args, "Warn")
}

func (l *LoggerEngineFile) Error(args ...interface{}) {
	filepath := l.Path + "/error.log"
	writeLog(filepath, args, "Error")
}

func (l *LoggerEngineFile) Fatal(args ...interface{}) {
	filepath := l.Path + "/fatal.log"
	writeLog(filepath, args, "Fatal")
}

func writeLog(filepath string, data []interface{}, label string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	var result map[string]interface{} = map[string]interface{}{}
	result["type"] = label
	result["time"] = time.Now()
	result["stack"] = string(debug.Stack())

	result["log"] = data
	buffer, err := json.Marshal(result)
	if err != nil {
		return err
	}
	buffer = append(buffer, 10)
	_, err = file.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}
