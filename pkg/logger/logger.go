package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	// InfoLogger is used for information messages
	InfoLogger *log.Logger
	// ErrorLogger is used for error messages
	ErrorLogger *log.Logger
)

// Init initializes the loggers
func Init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info logs an information message
func Info(format string, v ...interface{}) {
	if InfoLogger == nil {
		Init()
	}
	InfoLogger.Printf(format, v...)
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if ErrorLogger == nil {
		Init()
	}
	ErrorLogger.Printf(format, v...)
}

// LogRequest logs an HTTP request
func LogRequest(method, path string, status int, duration time.Duration) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s %s %d %v\n", timestamp, method, path, status, duration)
}
