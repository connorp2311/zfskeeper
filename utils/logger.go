package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	file   *os.File
	prefix string
}

func NewLogger(logFile string, prefix string) (*Logger, error) {
	var f *os.File
	if logFile != "" {
		logDir := filepath.Dir(logFile)
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			return nil, err
		}

		f, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}

	return &Logger{
		file:   f,
		prefix: prefix,
	}, nil
}

// Close closes the Logger and the underlying file.
func (l *Logger) Close() error {
	return l.file.Close()
}

//SetPrefix sets the prefix for the logger
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) Log(message string) {
	messageCon := fmt.Sprintf("\033[32m%s\033[0m \033[31m%s\033[0m %s", time.Now().Format("2006-01-02 15:04:05"), l.prefix, message)
	messageLog := fmt.Sprintf("%s %s %s", time.Now().Format("2006-01-02 15:04:05"), l.prefix, message)

	if l.file != nil {
		fmt.Fprintln(l.file, messageLog)
	}

	fmt.Println(messageCon)
}
