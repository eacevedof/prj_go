package components

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Logger writes timestamped log entries to dated files under the log directory.
// The directory is read from LOG_PATH env var; defaults to "logs".
type Logger struct {
	logDir string
}

// NewLogger returns a Logger bound to LOG_PATH (or "logs" if unset).
func NewLogger() *Logger {
	logDir := os.Getenv("LOG_PATH")
	if logDir == "" {
		logDir = "logs"
	}
	return &Logger{logDir: logDir}
}

func (l *Logger) LogError(module, message string, context map[string]any) {
	l.write("error", fmt.Sprintf("[%s] %s | %v", module, message, context))
}

func (l *Logger) LogInfo(module, message string) {
	l.write("info", fmt.Sprintf("[%s] %s", module, message))
}

func (l *Logger) LogDebug(module, message string, data any) {
	l.write("debug", fmt.Sprintf("[%s] %s | %v", module, message, data))
}

func (l *Logger) write(category, content string) {
	date := time.Now().Format("2006-01-02")
	ts := time.Now().Format("2006-01-02 15:04:05")
	path := filepath.Join(l.logDir, fmt.Sprintf("%s_%s.log", date, category))

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "[%s] %s\n", ts, content)
}
