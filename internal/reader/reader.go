// Package reader provides utilities for reading log input from various sources,
// including files, stdin, and compressed streams. It normalizes input into a
// channel of lines for downstream processing by the filter pipeline.
package reader

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// LineReader reads lines from an io.Reader and sends them to a channel.
type LineReader struct {
	// BufferSize controls the size of the internal scanner buffer.
	// Defaults to 1MB if not set.
	BufferSize int
}

// New returns a new LineReader with sensible defaults.
func New() *LineReader {
	return &LineReader{
		BufferSize: 1024 * 1024, // 1MB
	}
}

// ReadLines reads all lines from the given io.Reader and sends them to the
// returned channel. The channel is closed when the reader is exhausted or an
// error occurs. Errors are sent to the errCh channel.
func (r *LineReader) ReadLines(reader io.Reader) (<-chan string, <-chan error) {
	lines := make(chan string, 256)
	errs := make(chan error, 1)

	go func() {
		defer close(lines)
		defer close(errs)

		scanner := bufio.NewScanner(reader)
		buf := make([]byte, r.BufferSize)
		scanner.Buffer(buf, r.BufferSize)

		for scanner.Scan() {
			lines <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			errs <- fmt.Errorf("scanner error: %w", err)
		}
	}()

	return lines, errs
}

// OpenFile opens a file for reading, transparently decompressing gzip files
// based on their extension. The caller is responsible for closing the returned
// io.ReadCloser.
func OpenFile(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}

	if strings.HasSuffix(strings.ToLower(path), ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("creating gzip reader for %q: %w", path, err)
		}
		// Wrap both so closing the gzipReadCloser closes both the gzip reader
		// and the underlying file.
		return &gzipReadCloser{gr: gr, f: f}, nil
	}

	return f, nil
}

// Stdin returns os.Stdin as an io.ReadCloser for uniform handling.
func Stdin() io.ReadCloser {
	return os.Stdin
}

// IsCompressed reports whether the given file path appears to be a compressed
// file based on its extension. Currently only gzip (.gz) is supported.
func IsCompressed(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".gz")
}

// gzipReadCloser wraps a gzip.Reader and its underlying file so that both are
// closed when Close is called.
type gzipReadCloser struct {
	gr *gzip.Reader
	f  *os.File
}

func (g *gzipReadCloser) Read(p []byte) (int, error) {
	return g.gr.Read(p)
}

func (g *gzipReadCloser) Close() error {
	grErr := g.gr.Close()
	fErr := g.f.Close()
	if grErr != nil {
		return grErr
	}
	return fErr
}
