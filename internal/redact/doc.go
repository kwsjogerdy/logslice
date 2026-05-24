// Package redact provides a Redactor that scrubs sensitive data from log lines
// using user-supplied regular expressions.
//
// # Overview
//
// Create a Redactor with one or more regex patterns and an optional replacement
// placeholder. Each call to Apply or ApplyBytes scans the line and replaces
// every match with the placeholder (default: "[REDACTED]").
//
// # Example
//
//	r, err := redact.New([]string{`password=\S+`}, "")
//	if err != nil { /* handle */ }
//	clean := r.Apply("password=hunter2 user=bob")
//	// clean == "[REDACTED] user=bob"
package redact
