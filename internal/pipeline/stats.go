package pipeline

import "fmt"

// Stats records counters for a single pipeline run.
type Stats struct {
	// Read is the total number of lines read from the input source.
	Read int
	// Filtered is the number of lines dropped by the filter or level check.
	Filtered int
	// Written is the number of lines successfully written to the output.
	Written int
}

// String returns a human-readable summary of the pipeline run stats.
func (s Stats) String() string {
	return fmt.Sprintf("read=%d filtered=%d written=%d", s.Read, s.Filtered, s.Written)
}

// FilterRatio returns the fraction of read lines that were filtered out.
// Returns 0 if no lines were read.
func (s Stats) FilterRatio() float64 {
	if s.Read == 0 {
		return 0
	}
	return float64(s.Filtered) / float64(s.Read)
}
