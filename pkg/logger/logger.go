package logger

import (
    "fmt"
    "time"
)

// Константы уровней логирования
const (
    INFO  = "INFO"
    DEBUG = "DEBUG"
    ERROR = "ERROR"
)

// Logger структура логгера
type Logger struct{}

// New создает новый экземпляр логгера
func New() *Logger {
    return &Logger{}
}

// Info логирует информационные сообщения
func (l *Logger) Info(msg string) {
    fmt.Println(l.formatMessage(INFO, msg))
}

// Debug логирует отладочные сообщения
func (l *Logger) Debug(msg string) {
    fmt.Println(l.formatMessage(DEBUG, msg))
}

// Error логирует ошибки
func (l *Logger) Error(msg string, err error) {
    if err != nil {
        fmt.Println(l.formatMessage(ERROR, msg+" err - "+err.Error()))
    } else {
        fmt.Println(l.formatMessage(ERROR, msg))
    }
}

// formatMessage форматирует сообщение лога
func (l *Logger) formatMessage(level string, msg string) string {
    timeStr := time.Now().Format(time.RFC3339)
    return fmt.Sprintf(
        "[%s] %s, message - %s",
        level,
        timeStr,
        msg,
    )
}