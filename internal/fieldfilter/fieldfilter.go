// Package fieldfilter provides field-level filtering for structured log lines.
// It allows selecting or dropping specific fields from parsed key-value maps
// before they reach the formatter, enabling cleaner output for noisy logs.
package fieldfilter

import "strings"

// FieldFilter selects or drops fields from a parsed log record.
type FieldFilter struct {
	include map[string]struct{}
	exclude map[string]struct{}
}

// New creates a FieldFilter. includeFields limits output to only those keys
// (empty means keep all). excludeFields removes the named keys from output.
// Keys are matched case-insensitively.
func New(includeFields, excludeFields []string) *FieldFilter {
	ff := &FieldFilter{
		include: make(map[string]struct{}, len(includeFields)),
		exclude: make(map[string]struct{}, len(excludeFields)),
	}
	for _, f := range includeFields {
		ff.include[strings.ToLower(f)] = struct{}{}
	}
	for _, f := range excludeFields {
		ff.exclude[strings.ToLower(f)] = struct{}{}
	}
	return ff
}

// Enabled reports whether the filter will actually modify records.
func (ff *FieldFilter) Enabled() bool {
	return len(ff.include) > 0 || len(ff.exclude) > 0
}

// Apply returns a new map containing only the fields that pass the filter
// rules. The original map is not modified.
func (ff *FieldFilter) Apply(fields map[string]string) map[string]string {
	if !ff.Enabled() {
		return fields
	}
	out := make(map[string]string, len(fields))
	for k, v := range fields {
		lk := strings.ToLower(k)
		if len(ff.include) > 0 {
			if _, ok := ff.include[lk]; !ok {
				continue
			}
		}
		if _, ok := ff.exclude[lk]; ok {
			continue
		}
		out[k] = v
	}
	return out
}
