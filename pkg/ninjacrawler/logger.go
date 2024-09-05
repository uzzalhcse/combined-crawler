package ninjacrawler

import (
	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// logger is an interface for logging.
type logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Summary(format string, args ...interface{})
	Html(page playwright.Page, message string)
}

// defaultLogger is a default implementation of the logger interface using the standard log package.
type defaultLogger struct {
	logger         *log.Logger
	app            *Crawler
	gcpLogger      *log.Logger // GCP Summary logger
	gcpDebugLogger *log.Logger // GCP Debug logger
}

// newDefaultLogger creates a new instance of defaultLogger.
func newDefaultLogger(app *Crawler, siteName string) *defaultLogger {
	// Open a log file in append mode, create if it doesn't exist.

	currentDate := time.Now().Format("2006-01-02")
	directory := filepath.Join("storage", "logs", siteName)
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Construct the log file path.
	logFilePath := filepath.Join(directory, currentDate+"_application.log")

	// Open the log file in append mode, create if it doesn't exist.
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// Create a multi-writer that writes to both the file and the terminal.
	multiWriter := io.MultiWriter(file, os.Stdout)

	// Create the default logger
	dLogger := &defaultLogger{
		logger: log.New(multiWriter, "【"+app.Name+"】", log.LstdFlags),
		app:    app,
	}

	// Initialize GCP logger if requested
	if metadata.OnGCE() {
		dLogger.gcpLogger = getGCPLogger(app.Config, "ninjacrawler_log")
		dLogger.gcpDebugLogger = getGCPLogger(app.Config, "ninjacrawler_debug_log")
	}

	return dLogger

}

func getGCPLogger(config *configService, logID string) *log.Logger {

	os.Setenv("GCP_LOG_CREDENTIALS_PATH", "log-key.json")
	projectID := config.EnvString("PROJECT_ID", "lazuli-venturas-stg")

	client, err := logging.NewClient(context.Background(), projectID)
	if err != nil {
		panic(fmt.Sprintf("Failed to create GCP logging client: %v", err))
	}

	logger := client.Logger(logID)
	return logger.StandardLogger(logging.Debug)
}

// logWithGCP logs both to the local logger
func (l *defaultLogger) logWithGCP(level string, format string, args ...interface{}) {
	// Log to local logger
	l.logger.Printf(level+format, args...)

	// log to GCP
	if l.gcpLogger != nil {
		l.gcpLogger.Printf(format, args...)
	}
}
func (l *defaultLogger) Summary(format string, args ...interface{}) {
	l.logWithGCP("", format, args...)
}
func (l *defaultLogger) Debug(format string, args ...interface{}) {
	l.logger.Printf("DEBUG: "+format, args...)
	// log to GCP
	if l.gcpDebugLogger != nil {
		l.gcpDebugLogger.Printf(format, args...)
	}
}
func (l *defaultLogger) Info(format string, args ...interface{}) {
	l.logger.Printf("✔ "+format, args...)
}
func (l *defaultLogger) Warn(format string, args ...interface{}) {
	l.logger.Printf("⚠️ "+format, args...)
}

func (l *defaultLogger) Error(format string, args ...interface{}) {
	l.logger.Printf("🛑 ERROR: "+format, args...)
}

func (l *defaultLogger) Fatal(format string, args ...interface{}) {
	l.logger.Fatalf("🚨 FATAL: "+format, args...)
}

func (l *defaultLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *defaultLogger) Html(html, url, msg string) {
	if l.app.IsValidPage(url) {
		l.Error("Html Error: %v", msg)
		err := l.app.writePageContentToFile(html, url, msg)
		if err != nil {
			l.logger.Printf("⚛️ HTML: %v", err)
		}
	}
}
