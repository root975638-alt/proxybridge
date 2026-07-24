// Package logging provides structured logging for ProxyBridge.
package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level represents log severity
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Logger holds logging configuration
type Logger struct {
	level       Level
	jsonOutput  bool
	mu          sync.Mutex
	logFile     *os.File
	logDir      string
	maxLogSize  int64
	maxLogAge   int
}

var (
	globalLogger *Logger
	once         sync.Once
)

// Init initializes the global logger
func Init(level string, jsonOutput bool) {
	once.Do(func() {
		lvl := InfoLevel
		switch strings.ToLower(level) {
		case "debug":
			lvl = DebugLevel
		case "info":
			lvl = InfoLevel
		case "warn", "warning":
			lvl = WarnLevel
		case "error":
			lvl = ErrorLevel
		case "fatal":
			lvl = FatalLevel
		}

		logDir, err := getLogDirectory()
		if err != nil {
			logDir = ""
		}

		globalLogger = &Logger{
			level:       lvl,
			jsonOutput:  jsonOutput,
			logDir:      logDir,
			maxLogSize:  10 * 1024 * 1024, // 10MB
			maxLogAge:   30, // days
		}
	})
}

// GetLogger returns the global logger
func GetLogger() *Logger {
	if globalLogger == nil {
		Init("info", false)
	}
	return globalLogger
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetJSONOutput sets JSON output mode
func (l *Logger) SetJSONOutput(jsonOutput bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.jsonOutput = jsonOutput
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.log(DebugLevel, msg, keysAndValues...)
}

// Info logs an info message
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.log(InfoLevel, msg, keysAndValues...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.log(WarnLevel, msg, keysAndValues...)
}

// Error logs an error message
func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.log(ErrorLevel, msg, keysAndValues...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log(FatalLevel, msg, keysAndValues...)
	os.Exit(1)
}

// log writes a log entry
func (l *Logger) log(level Level, msg string, keysAndValues ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	// Skip if we're using JSON output and at debug level (too noisy)
	if l.jsonOutput && level == DebugLevel {
		return
	}

	// Build log entry
	entry := l.buildEntry(level, msg, keysAndValues...)

	// Write to stderr
	fmt.Fprintln(os.Stderr, entry)

	// Write to log file if available
	if l.logFile != nil {
		fmt.Fprintln(l.logFile, entry)
	}
}

// buildEntry builds a log entry
func (l *Logger) buildEntry(level Level, msg string, keysAndValues ...interface{}) string {
	if l.jsonOutput {
		return l.buildJSONEntry(level, msg, keysAndValues...)
	}
	return l.buildTextEntry(level, msg, keysAndValues...)
}

// buildJSONEntry builds a JSON log entry
func (l *Logger) buildJSONEntry(level Level, msg string, keysAndValues ...interface{}) string {
	fields := l.buildFields(msg, keysAndValues...)

	entry := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"level":     level.String(),
		"message":   msg,
		"fields":    fields,
	}

	if len(keysAndValues) > 0 {
		// Extract additional context
		if len(keysAndValues)%2 == 0 {
			for i := 0; i < len(keysAndValues); i += 2 {
				if key, ok := keysAndValues[i].(string); ok {
					entry[key] = keysAndValues[i+1]
				}
			}
		}
	}

	data, _ := json.Marshal(entry)
	return string(data)
}

// buildTextEntry builds a text log entry
func (l *Logger) buildTextEntry(level Level, msg string, keysAndValues ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := level.String()

	// Build additional fields
	var extra []string
	if len(keysAndValues) > 0 {
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				if key, ok := keysAndValues[i].(string); ok {
					extra = append(extra, fmt.Sprintf("%s=%v", key, keysAndValues[i+1]))
				}
			}
		}
	}

	var result string
	if len(extra) > 0 {
		result = fmt.Sprintf("[%s] %s %s | %s", levelStr, timestamp, msg, strings.Join(extra, " "))
	} else {
		result = fmt.Sprintf("[%s] %s %s", levelStr, timestamp, msg)
	}

	return result
}

// buildFields extracts log fields
func (l *Logger) buildFields(msg string, keysAndValues ...interface{}) map[string]interface{} {
	fields := make(map[string]interface{})

	if len(keysAndValues) > 0 {
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				if key, ok := keysAndValues[i].(string); ok {
					fields[key] = keysAndValues[i+1]
				}
			}
		}
	}

	// Add caller info
	if _, file, line, ok := runtime.Caller(3); ok {
		fields["caller"] = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	return fields
}

// String returns string representation of level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Debug logs a debug message
func Debug(msg string, keysAndValues ...interface{}) {
	GetLogger().Debug(msg, keysAndValues...)
}

// Info logs an info message
func Info(msg string, keysAndValues ...interface{}) {
	GetLogger().Info(msg, keysAndValues...)
}

// Warn logs a warning message
func Warn(msg string, keysAndValues ...interface{}) {
	GetLogger().Warn(msg, keysAndValues...)
}

// Error logs an error message
func Error(msg string, keysAndValues ...interface{}) {
	GetLogger().Error(msg, keysAndValues...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, keysAndValues ...interface{}) {
	GetLogger().Fatal(msg, keysAndValues...)
	os.Exit(1)
}

// getLogDirectory returns the log directory
func getLogDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	logDir := filepath.Join(home, ".config", "proxybridge", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", err
	}

	return logDir, nil
}

// RotateLogs rotates log files
func RotateLogs() error {
	l := GetLogger()
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logDir == "" {
		return fmt.Errorf("log directory not set")
	}

	// Implementation would rotate logs based on size/age
	// For now, just a placeholder

	return nil
}

// Close closes the logger
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logFile != nil {
		l.logFile.Close()
		l.logFile = nil
	}

	return nil
}
