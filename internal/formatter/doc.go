// Package formatter provides log line rendering utilities for logslice.
//
// It is responsible for mapping severity level strings (e.g. "info", "ERROR")
// to a canonical Level value, and for optionally applying ANSI color codes to
// log lines based on their severity.
//
// Usage:
//
//	f := formatter.New(true) // true = enable colorized output
//	formatted := f.Format(line, levelString)
//
// When colorization is disabled the original line is returned unchanged,
// making the formatter safe to use when writing to non-TTY outputs.
package formatter
