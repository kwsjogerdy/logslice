// Package output provides utilities for writing filtered log lines
// to various destinations such as stdout or files.
package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Writer wraps an io.Writer with buffered output and optional colorization.
type Writer struct {
	w       *bufio.Writer
	colorize bool
}

// ANSI color codes for highlighted output.
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// New creates a Writer that writes to the given io.Writer.
// If colorize is true, matched lines are highlighted using ANSI codes.
func New(w io.Writer, colorize bool) *Writer {
	return &Writer{
		w:        bufio.NewWriter(w),
		colorize: colorize,
	}
}

// Stdout returns a Writer that writes to os.Stdout.
func Stdout(colorize bool) *Writer {
	return New(os.Stdout, colorize)
}

// OpenFile returns a Writer that writes to the named file, creating or
// truncating it as needed. The caller is responsible for closing the file.
func OpenFile(path string, colorize bool) (*Writer, *os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("output: open file %q: %w", path, err)
	}
	return New(f, colorize), f, nil
}

// WriteLine writes a single log line followed by a newline character.
// If colorize is enabled the line is wrapped in green ANSI escape codes.
func (w *Writer) WriteLine(line string) error {
	var err error
	if w.colorize {
		_, err = fmt.Fprintf(w.w, "%s%s%s\n", colorGreen, line, colorReset)
	} else {
		_, err = fmt.Fprintln(w.w, line)
	}
	return err
}

// Flush flushes any buffered data to the underlying writer.
func (w *Writer) Flush() error {
	return w.w.Flush()
}
