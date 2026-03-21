package logger

import (
    "fmt"
    "time"
)

type SimpleLogger struct {
    level string
}

func NewSimpleLogger() *SimpleLogger {
    return &SimpleLogger{level: "INFO"}
}

func (l *SimpleLogger) Info(msg string) {
    fmt.Printf("[%s] INFO: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}

func (l *SimpleLogger) Debug(msg string) {
    fmt.Printf("[%s] DEBUG: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}

func (l *SimpleLogger) Error(msg string) {
    fmt.Printf("[%s] ERROR: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}