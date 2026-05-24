// Package highlight implements substring search and annotation for log lines.
//
// # Overview
//
// A Highlighter is created with a set of search terms and a case-sensitivity
// flag. Calling Scan on a log line returns a Result that lists every matched
// byte range inside the line. The Annotate helper then wraps those ranges with
// caller-supplied prefix/suffix strings, making it straightforward to inject
// ANSI colour codes or any other decoration before the line reaches the output
// writer.
//
// # Usage
//
//	h := highlight.New([]string{"error", "warn"}, true)
//	result := h.Scan(line)
//	annotated := highlight.Annotate(result, "\033[1;31m", "\033[0m")
//
// The package is intentionally free of I/O concerns; callers decide how and
// whether to render the annotated output.
package highlight
