// Package logger is a minimal structured logger used across the monorepo's
// independent Go modules (shopping-svc, estargz, ipfs) to prove that a
// single shared-lib module can be required by all of them at once.
package logger

import (
	"fmt"
	"time"
)

// Level identifies a log severity.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Logger emits timestamped, tagged log lines.
type Logger struct {
	component string
}

// New returns a Logger tagged with the given component name.
func New(component string) *Logger {
	return &Logger{component: component}
}

func (l *Logger) log(level Level, msg string, args ...any) {
	fmt.Printf("%s [%s] (%s) %s\n",
		time.Now().UTC().Format(time.RFC3339),
		level,
		l.component,
		fmt.Sprintf(msg, args...),
	)
}

func (l *Logger) Info(msg string, args ...any)  { l.log(LevelInfo, msg, args...) }
func (l *Logger) Warn(msg string, args ...any)  { l.log(LevelWarn, msg, args...) }
func (l *Logger) Error(msg string, args ...any) { l.log(LevelError, msg, args...) }
