package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

var (
	logger Logger
)

type Logger interface {
	Log(format string, args ...interface{})
}

// SetLogger sets the debug log file
func SetLogger(l Logger) {
	logger = l
}

// Enabled returns whether debug logging is enabled
func Enabled() bool {
	return logger != nil
}

// Logf logs debug messages to the debug file if available
func Logf(format string, args ...interface{}) {
	if logger != nil {
		// get file and line number
		_, file, line, _ := runtime.Caller(1)

		baseFile := filepath.Base(file)
		timestamp := time.Now().Format("2006-01-02 15:04:05.000")
		message := fmt.Sprintf(format, args...)
		logger.Log(fmt.Sprintf("[%s] %s:%d DEBUG: %s", timestamp, baseFile, line, message))
	}
}
