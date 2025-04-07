package notification

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ecg_tool/pkg/model"
)

// Logger handles logging of heart conditions
type Logger struct {
	fileLogger    *log.Logger
	consoleLogger *log.Logger
	mu            sync.Mutex
	logFile       *os.File
}

// NewLogger creates a new logger that writes to both console and file
func NewLogger(logDirPath string) (*Logger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Create log file with current timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(logDirPath, fmt.Sprintf("ecg_log_%s.log", timestamp))

	logFile, err := os.Create(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %v", err)
	}

	fileLogger := log.New(logFile, "", log.LstdFlags)
	consoleLogger := log.New(os.Stdout, "", log.LstdFlags)

	return &Logger{
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
		logFile:       logFile,
	}, nil
}

// LogCondition logs a heart condition
func (l *Logger) LogCondition(condition *model.HeartCondition) {
	if condition == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	logMsg := fmt.Sprintf("ALERT: %s detected - Severity: %s, Heart rate: %d BPM at %s - %s",
		condition.Type,
		condition.Severity,
		condition.HeartRate,
		condition.Timestamp.Format("15:04:05"),
		condition.Description)

	// Log to file
	l.fileLogger.Println(logMsg)

	// Log to console
	l.consoleLogger.Println(logMsg)
}

// LogNormalReading logs a normal heart reading
func (l *Logger) LogNormalReading(data model.ECGData) {
	l.mu.Lock()
	defer l.mu.Unlock()

	logMsg := fmt.Sprintf("NORMAL: Heart rate: %d BPM, RR Interval: %.2fs at %s",
		data.HeartRate,
		data.RRInterval,
		data.Timestamp.Format("15:04:05"))

	// Log only to file to avoid console clutter
	l.fileLogger.Println(logMsg)
}

// Close closes the log file
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
