// Package highlight provides term-matching and substring highlighting
// for log lines, marking matched regions so the output layer can
// render them with colour or other emphasis.
package highlight

import (
	"strings"
)

// Match describes a single matched region within a line.
type Match struct {
	Start int
	End   int
}

// Result holds the original line together with every match found
// inside it. When Matches is empty the line contained no hits.
type Result struct {
	Line    string
	Matches []Match
}

// Highlighter searches lines for one or more terms.
type Highlighter struct {
	terms     []string
	caseBlind bool
}

// New returns a Highlighter that searches for the given terms.
// If caseInsensitive is true both the line and each term are
// lower-cased before comparison (original offsets are still
// reported against the original line).
func New(terms []string, caseInsensitive bool) *Highlighter {
	norm := make([]string, len(terms))
	for i, t := range terms {
		if caseInsensitive {
			norm[i] = strings.ToLower(t)
		} else {
			norm[i] = t
		}
	}
	return &Highlighter{terms: norm, caseBlind: caseInsensitive}
}

// Scan searches line for every occurrence of every term and returns
// a Result. Overlapping matches are preserved as-is.
func (h *Highlighter) Scan(line string) Result {
	search := line
	if h.caseBlind {
		search = strings.ToLower(line)
	}

	var matches []Match
	for _, term := range h.terms {
		if term == "" {
			continue
		}
		offset := 0
		for {
			idx := strings.Index(search[offset:], term)
			if idx < 0 {
				break
			}
			start := offset + idx
			end := start + len(term)
			matches = append(matches, Match{Start: start, End: end})
			offset = end
		}
	}
	return Result{Line: line, Matches: matches}
}

// Annotate wraps every matched region in line with the supplied prefix
// and suffix strings (e.g. ANSI escape codes) and returns the resulting
// string. If there are no matches the original line is returned unchanged.
func Annotate(r Result, prefix, suffix string) string {
	if len(r.Matches) == 0 {
		return r.Line
	}
	var sb strings.Builder
	prev := 0
	for _, m := range r.Matches {
		sb.WriteString(r.Line[prev:m.Start])
		sb.WriteString(prefix)
		sb.WriteString(r.Line[m.Start:m.End])
		sb.WriteString(suffix)
		prev = m.End
	}
	sb.WriteString(r.Line[prev:])
	return sb.String()
}
