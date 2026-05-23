package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Format represents the detected log format.
type Format int

const (
	FormatPlain  Format = iota
	FormatJSON
	FormatLogfmt
)

// Entry represents a parsed log line with optional structured fields.
type Entry struct {
	Raw       string
	Format    Format
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]string
}

// Parse attempts to parse a raw log line into an Entry.
// It tries JSON first, then logfmt, then falls back to plain text.
func Parse(line string) Entry {
	entry := Entry{Raw: line}

	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return entry
	}

	if trimmed[0] == '{' {
		if e, ok := parseJSON(trimmed); ok {
			e.Raw = line
			return e
		}
	}

	if e, ok := parseLogfmt(trimmed); ok {
		e.Raw = line
		return e
	}

	entry.Format = FormatPlain
	entry.Message = trimmed
	return entry
}

func parseJSON(line string) (Entry, bool) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return Entry{}, false
	}

	entry := Entry{
		Format: FormatJSON,
		Fields: make(map[string]string),
	}

	for k, v := range raw {
		s := toString(v)
		switch strings.ToLower(k) {
		case "msg", "message":
			entry.Message = s
		case "level", "severity":
			entry.Level = strings.ToUpper(s)
		case "time", "timestamp", "ts":
			entry.Timestamp, _ = time.Parse(time.RFC3339, s)
		default:
			entry.Fields[k] = s
		}
	}
	return entry, true
}

func parseLogfmt(line string) (Entry, bool) {
	fields := make(map[string]string)
	parts := strings.Fields(line)
	found := false
	for _, part := range parts {
		idx := strings.IndexByte(part, '=')
		if idx < 1 {
			continue
		}
		found = true
		key := part[:idx]
		val := strings.Trim(part[idx+1:], `"`)
		fields[key] = val
	}
	if !found {
		return Entry{}, false
	}

	entry := Entry{Format: FormatLogfmt, Fields: fields}
	if v, ok := fields["msg"]; ok {
		entry.Message = v
		delete(fields, "msg")
	} else if v, ok := fields["message"]; ok {
		entry.Message = v
		delete(fields, "message")
	}
	if v, ok := fields["level"]; ok {
		entry.Level = strings.ToUpper(v)
		delete(fields, "level")
	}
	return entry, true
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strings.TrimRight(fmt.Sprintf("%g", t), ".")
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
