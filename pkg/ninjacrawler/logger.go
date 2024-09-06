package ninjacrawler

import (
	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"go.uber.org/zap"
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

// defaultLogger is a default implementation of the logger interface using Zap.
type defaultLogger struct {
	logger         *zap.SugaredLogger
	app            *Crawler
	gcpLogger      *logging.Logger // GCP Summary logger
	gcpDebugLogger *logging.Logger // GCP Debug logger
	siteName       string
}

// newDefaultLogger creates a new instance of defaultLogger.
func newDefaultLogger(app *Crawler, siteName string) *defaultLogger {
	// Create the log directory
	logFileName := getLogFileName(siteName)

	// Setup Zap logger to write to both file and console, logging only the message
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"stdout",
		logFileName,
	}
	cfg.EncoderConfig.EncodeTime = nil   // Disable timestamp in local logs
	cfg.EncoderConfig.TimeKey = ""       // No time key in local logs
	cfg.EncoderConfig.LevelKey = ""      // No log level in local logs
	cfg.EncoderConfig.CallerKey = ""     // No caller information in local logs
	cfg.EncoderConfig.MessageKey = "msg" // Keep only the message

	// Build the logger without caller and time information
	zapLogger, err := cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to create Zap logger: %v", err))
	}

	sugarLogger := zapLogger.Sugar()

	dLogger := &defaultLogger{
		logger:   sugarLogger,
		app:      app,
		siteName: siteName,
	}

	// Initialize GCP logger if requested
	if metadata.OnGCE() {
		dLogger.gcpLogger = getGCPLogger(app.Config, "ninjacrawler_summary_log")
		dLogger.gcpDebugLogger = getGCPLogger(app.Config, "ninjacrawler_dev_log")
	}

	return dLogger
}
func getLogFileName(siteName string) string {
	currentDate := time.Now().Format("2006-01-02")
	directory := filepath.Join("storage", "logs", siteName)
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Construct the log file path.
	logFilePath := filepath.Join(directory, currentDate+"_application.log")
	return logFilePath
}

func getGCPLogger(config *configService, logID string) *logging.Logger {
	os.Setenv("GCP_LOG_CREDENTIALS_PATH", "log-key.json")
	projectID := config.EnvString("PROJECT_ID", "lazuli-venturas-stg")

	client, err := logging.NewClient(context.Background(), projectID)
	if err != nil {
		panic(fmt.Sprintf("Failed to create GCP logging client: %v", err))
	}

	logger := client.Logger(logID)
	return logger
}

// logWithGCP logs both to the local logger and GCP.
func (l *defaultLogger) logWithGCP(level string, msg string, args ...interface{}) {
	ts := time.Now().Format("2006-01-02 15:04:05")

	// Log to local logger (console and file) - only the message
	l.logger.Infof(msg, args...)

	// Log to GCP
	if l.gcpLogger != nil && level != "summary" {
		l.gcpLogger.Log(logging.Entry{
			Payload: map[string]interface{}{
				"level":     "info",
				"caller":    "", // Caller information is logged automatically in the local logger.
				"ts":        ts,
				"site_name": l.siteName,
				"msg":       fmt.Sprintf(msg, args...),
			},
			Severity: logging.Default,
		})
	}

	// Log debug to GCP Debug logger
	if l.gcpDebugLogger != nil && level == "debug" {
		l.gcpDebugLogger.Log(logging.Entry{
			Payload: map[string]interface{}{
				"level":     "error",
				"caller":    "", // Caller information is logged automatically in the local logger.
				"ts":        ts,
				"site_name": l.siteName,
				"msg":       fmt.Sprintf(msg, args...),
			},
			Severity: logging.Debug,
		})
	}
}

func (l *defaultLogger) Summary(format string, args ...interface{}) {
	l.logWithGCP("summary", format, args...)
}

func (l *defaultLogger) Debug(format string, args ...interface{}) {
	l.logWithGCP("debug", format, args...)
}

func (l *defaultLogger) Info(format string, args ...interface{}) {
	l.logWithGCP("info", "✔ "+format, args...)
}
func (l *defaultLogger) Warn(format string, args ...interface{}) {
	l.logWithGCP("warn", "⚠️ "+format, args...)
}

func (l *defaultLogger) Error(format string, args ...interface{}) {
	l.logWithGCP("error", "🛑 "+format, args...)
}

func (l *defaultLogger) Fatal(format string, args ...interface{}) {
	l.logWithGCP("fatal", "🚨 "+format, args...)
	os.Exit(1)
}

func (l *defaultLogger) Printf(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *defaultLogger) Html(html, url, msg string) {
	if l.app.IsValidPage(url) {
		l.Error("Html Error: %v", msg)
		err := l.app.writePageContentToFile(html, url, msg)
		if err != nil {
			l.logger.Infof("⚛️ HTML: %v", err)
		}
	}
}
