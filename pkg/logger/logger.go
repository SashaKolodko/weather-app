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

// Структура логгера
type Logger struct{}

// Конструктор логгера
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

// Error логирует ошибки - принимает string и error
func (l *Logger) Error(msg string, err error) {
    fmt.Println(l.formatMessage(ERROR, msg+" err - "+err.Error()))
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