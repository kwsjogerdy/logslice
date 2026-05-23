// Package parser provides log line parsing for logslice.
//
// It supports three log formats:
//
//   - JSON: lines beginning with '{' are parsed as JSON objects.
//     Common fields such as "msg"/"message", "level"/"severity", and
//     "time"/"timestamp"/"ts" are promoted to first-class Entry fields.
//
//   - Logfmt: lines containing key=value pairs (optionally quoted) are
//     parsed as logfmt. The "level" and "msg"/"message" keys are promoted
//     to first-class Entry fields.
//
//   - Plain: any line that does not match the above formats is treated as
//     unstructured plain text. The entire trimmed line becomes the Message.
//
// The original raw line is always preserved in Entry.Raw regardless of
// the detected format, ensuring lossless round-trips through the pipeline.
//
// Usage:
//
//	entry := parser.Parse(line)
//	fmt.Println(entry.Level, entry.Message)
package parser
