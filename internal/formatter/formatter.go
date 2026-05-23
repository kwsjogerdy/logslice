package formatter

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Level represents a log severity level.
type Level int

const (
	LevelUnknown Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// ParseLevel maps common severity strings to a Level.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug", "dbg", "trace":
		return LevelDebug
	case "info", "information":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error", "err":
		return LevelError
	case "fatal", "critical", "panic":
		return LevelFatal
	default:
		return LevelUnknown
	}
}

// Formatter controls how log lines are rendered.
type Formatter struct {
	Colorize bool
}

// New returns a new Formatter.
func New(colorize bool) *Formatter {
	return &Formatter{Colorize: colorize}
}

// Format renders a log line, optionally applying color based on severity.
func (f *Formatter) Format(line string, levelStr string) string {
	if !f.Colorize {
		return line
	}

	lvl := ParseLevel(levelStr)
	var fn func(string, ...interface{}) string

	switch lvl {
	case LevelDebug:
		fn = color.New(color.FgCyan).SprintfFunc()
	case LevelInfo:
		fn = color.New(color.FgGreen).SprintfFunc()
	case LevelWarn:
		fn = color.New(color.FgYellow).SprintfFunc()
	case LevelError:
		fn = color.New(color.FgRed).SprintfFunc()
	case LevelFatal:
		fn = color.New(color.FgRed, color.Bold).SprintfFunc()
	default:
		return line
	}

	return fmt.Sprintf("%s", fn("%s", line))
}
