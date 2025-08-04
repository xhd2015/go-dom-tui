package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	logger Logger
)

type Logger interface {
	Logf(format string, args ...interface{})
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
		logger.Logf("["+timestamp+"] %s:%d DEBUG: %s", baseFile, line, message)
	}
}

func NewFileLogger(file *os.File) *fileLogger {
	return &fileLogger{file: file}
}

type fileLogger struct {
	file *os.File
}

func (l *fileLogger) Logf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.file.Write([]byte(msg))
	if !strings.HasSuffix(msg, "\n") {
		l.file.Write([]byte("\n"))
	}
}
