package client

import (
	"fmt"
	"log"
	"os"
)

// SimpleLogger простая реализация Logger используя стандартный log пакет
type SimpleLogger struct {
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	level    LogLevel
}

// LogLevel уровень логирования
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// NewSimpleLogger создает простой логгер
func NewSimpleLogger(level LogLevel) *SimpleLogger {
	flags := log.LstdFlags | log.Lshortfile
	return &SimpleLogger{
		debugLog: log.New(os.Stdout, "[DEBUG] ", flags),
		infoLog:  log.New(os.Stdout, "[INFO] ", flags),
		warnLog:  log.New(os.Stderr, "[WARN] ", flags),
		errorLog: log.New(os.Stderr, "[ERROR] ", flags),
		level:    level,
	}
}

func (s *SimpleLogger) Debug(msg string, fields ...Field) {
	if s.level <= LogLevelDebug {
		s.log(s.debugLog, msg, fields...)
	}
}

func (s *SimpleLogger) Info(msg string, fields ...Field) {
	if s.level <= LogLevelInfo {
		s.log(s.infoLog, msg, fields...)
	}
}

func (s *SimpleLogger) Warn(msg string, fields ...Field) {
	if s.level <= LogLevelWarn {
		s.log(s.warnLog, msg, fields...)
	}
}

func (s *SimpleLogger) Error(msg string, fields ...Field) {
	if s.level <= LogLevelError {
		s.log(s.errorLog, msg, fields...)
	}
}

func (s *SimpleLogger) log(logger *log.Logger, msg string, fields ...Field) {
	if len(fields) == 0 {
		logger.Println(msg)
		return
	}

	fieldsStr := ""
	for i, field := range fields {
		if i > 0 {
			fieldsStr += ", "
		}
		fieldsStr += fmt.Sprintf("%s=%v", field.Key, field.Value)
	}
	logger.Printf("%s [%s]", msg, fieldsStr)
}

