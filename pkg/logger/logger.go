// Package logger provides structured leveled logging for the application.
// It writes logs to both console and a timestamped log file.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/antonlearn/go-final-project/pkg/settings"
)

var (
	// Leveled loggers
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

// SetupLogger initializes the leveled loggers and creates a timestamped log file.
// It returns the file handle for proper closing in main.
func SetupLogger() (*os.File, error) {
	// Generate timestamped log filename
	logFileName := generationLocalFileName(".log")

	file, err := os.Create(logFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %s: %w", logFileName, err)
	}

	// Write logs to file/stdout/file and stdout
	var multiWriter io.Writer
	switch settings.LogFileStdout {
	case "File":
		multiWriter = io.MultiWriter(file)
	case "Stdout":
		multiWriter = io.MultiWriter(os.Stdout)
	default:
		multiWriter = io.MultiWriter(file, os.Stdout)
	}

	// Initialize leveled loggers with timestamp and file information
	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(multiWriter, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return file, nil
}

// generationLocalFileName generates a unique log filename with current UTC timestamp
func generationLocalFileName(ext string) string {
	timeStamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	return "app_" + timeStamp + ext
}

// Info logs a message at INFO level
func Info(msg string) {
	InfoLogger.Println(msg)
}

// Infof logs a formatted message at INFO level
func Infof(format string, v ...any) {
	InfoLogger.Printf(format, v...)
}

// Warn logs a message at WARN level
func Warn(msg string) {
	WarnLogger.Println(msg)
}

// Warnf logs a formatted message at WARN level
func Warnf(format string, v ...any) {
	WarnLogger.Printf(format, v...)
}

// Error logs a message at ERROR level
func Error(msg string) {
	ErrorLogger.Println(msg)
}

// Errorf logs a formatted message at ERROR level
func Errorf(format string, v ...any) {
	ErrorLogger.Printf(format, v...)
}
