package utils

import (
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

type Logger struct {
	lastFlushed time.Time
	writer      io.Writer
	mux         sync.Mutex
}

func NewLogger(writer io.Writer) *Logger {
	if writer == nil {
		panic("Cannot pass nil writer to logger.")
	}
	return &Logger{
		writer: writer,
		mux:    sync.Mutex{},
	}
}

// logs a job message to writer.
func (l *Logger) LogJob(name, message string) (int, error) {
	day := time.Now().Format("2006-01-02 15:04:00")
	l.mux.Lock()
	defer l.mux.Unlock()
	n, err := l.writer.Write([]byte(fmt.Sprintf("%s: [%s] %s\n", day, name, message)))
	if err != nil {
		slog.Error("Failed to write to log file", "Error", err)
		return 0, err
	}
	fmt.Printf("%s: %s logged successfully\n", day, name)
	return n, nil
}

// logs an error message to writer.
func (l *Logger) LogError(err error) (int, error) {
	if err == nil {
		return 0, nil
	}
	day := time.Now().Format("2006-01-02 15:04:00")

	// Get the error message
	message := err.Error()

	l.mux.Lock()
	defer l.mux.Unlock()
	n, err := l.writer.Write([]byte(fmt.Sprintf("%s: %s\n", day, message)))
	if err != nil {
		slog.Error("Failed to write to log file", "Error", err)
		return 0, err
	}
	// Log to console
	slog.Error("Error:", "Time", day, "Message", message)
	return n, nil
}
