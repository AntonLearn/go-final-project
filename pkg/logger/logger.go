// Package logger provides a structured, leveled logging implementation that wraps
// standard library log primitives to support multi-stream outputs (file, stdout, or both).
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Interface defines the decoupled logging contract enforced across various application
// architecture layers.
type Interface interface {
	Info(v ...any)
	Infof(format string, v ...any)
	Warn(v ...any)
	Warnf(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
}

// Logger implements the Interface contract by encapsulating isolated standard log
// instances mapped to distinct operational severity levels.
type Logger struct {
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	file     *os.File
}

// New initializes and returns a configured Logger instance. Depending on the specified
// output mode ("File", "Stdout", "Both"), it provisions underlying writer streams and
// handles unique, timestamped log file allocations.
func New(mode string) (*Logger, error) {
	var file *os.File
	var err error
	var writer io.Writer

	// Provision a local log file resource if the requested execution mode demands persistence.
	if mode == "File" || mode == "Both" || mode == "" {
		logFileName := generateLocalFileName(".log")
		file, err = os.Create(logFileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file %s: %w", logFileName, err)
		}
	}

	// Calibrate downstream destination streams based on operational topology preferences.
	switch mode {
	case "File":
		writer = file
	case "Stdout":
		writer = os.Stdout
	default: // Handles "Both" or unassigned fallback configurations
		if file != nil {
			writer = io.MultiWriter(file, os.Stdout)
		} else {
			writer = os.Stdout
		}
	}

	// Capture date, precise time coordinates, and the immediate short file name descriptor.
	flags := log.Ldate | log.Ltime | log.Lshortfile

	return &Logger{
		infoLog:  log.New(writer, "INFO:  ", flags),
		warnLog:  log.New(writer, "WARN:  ", flags),
		errorLog: log.New(writer, "ERROR: ", flags),
		file:     file,
	}, nil
}

// Close gracefully releases the underlying log file descriptor if it was allocated
// during initialization.
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Info logs messages to the informational stream using default line formatting.
func (l *Logger) Info(v ...any) { l.infoLog.Println(v...) }

// Infof logs formatted messages to the informational stream using standard template evaluation.
func (l *Logger) Infof(format string, v ...any) { l.infoLog.Printf(format, v...) }

// Warn logs warning alerts to the tracking stream using default line formatting.
func (l *Logger) Warn(v ...any) { l.warnLog.Println(v...) }

// Warnf logs formatted warning alerts to the tracking stream using standard template evaluation.
func (l *Logger) Warnf(format string, v ...any) { l.warnLog.Printf(format, v...) }

// Error logs failure messages to the system error stream using default line formatting.
func (l *Logger) Error(v ...any) { l.errorLog.Println(v...) }

// Errorf logs formatted failure messages to the system error stream using standard template evaluation.
func (l *Logger) Errorf(format string, v ...any) { l.errorLog.Printf(format, v...) }

// generateLocalFileName crafts a unique file naming constraint incorporating an active
// UTC execution timestamp to avoid log collision risks.
func generateLocalFileName(ext string) string {
	timeStamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	return "app_" + timeStamp + ext
}
