// Package logger
package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/antonlearn/go-final-project/pkg/config"
)

func SetupLogger() (*os.File, error) {
	file, err := os.Create(generationLocalFileName(".log"))
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %s: %w", file.Name(), err)
	}
	// Creating new logger
	config.Config.Logger = log.New(file, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return file, nil
}

func generationLocalFileName(ext string) string {
	timeStamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	fileName := "app_" + timeStamp + ext
	return fileName
}
