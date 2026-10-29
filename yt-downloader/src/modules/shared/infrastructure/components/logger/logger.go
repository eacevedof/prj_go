package logger

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"yt-downloader-go/src/modules/shared/infrastructure/repositories/configuration"
)

var writeLogMutex sync.Mutex

// Logger escribe en <logsPath>/YYYY-MM-DD-<nombre>.log
type Logger struct{}

func GetInstance() *Logger {
	return &Logger{}
}

func (l *Logger) WriteLog(filePath string, content string, logsPath string) {
	if logsPath == "" {
		logsPath = configuration.GetInstance().GetLogPath()
	}
	if logsPath == "" {
		return
	}
	if err := os.MkdirAll(logsPath, 0o777); err != nil {
		return
	}
	logsFolderPath, err := filepath.Abs(logsPath)
	if err != nil {
		return
	}
	logFilePath := filepath.Join(logsFolderPath, filePath)
	finalDir := filepath.Dir(logFilePath)
	_ = os.MkdirAll(finalDir, 0o777)

	ext := strings.TrimPrefix(filepath.Ext(logFilePath), ".")
	fileName := strings.TrimSuffix(filepath.Base(logFilePath), filepath.Ext(logFilePath))
	if ext == "" {
		ext = "log"
	}

	// la app trabaja en UTC
	now := time.Now().UTC()
	logFilePath = filepath.Join(finalDir, now.Format("2006-01-02")+"-"+fileName+"."+ext)
	content = strings.TrimSpace(content)

	writeLogMutex.Lock()
	defer writeLogMutex.Unlock()
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return
	}
	defer logFile.Close()
	_, _ = logFile.WriteString("\n[" + now.Format("2006-01-02 15:04:05") + "]\n" + content)
}
