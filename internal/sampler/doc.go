// Package sampler provides probabilistic log-line sampling for logslice.
//
// When processing high-volume log streams it is often useful to inspect only
// a representative fraction of the data.  The Sampler type wraps a fast
// pseudo-random number generator and exposes a single Keep method that callers
// invoke once per log line.
//
// Rates:
//
//	1.0  — keep every line (default, no overhead beyond the rate check)
//	0.5  — keep roughly half of all lines
//	0.0  — discard every line
//
// Thread safety: Keep is NOT safe for concurrent use from multiple goroutines.
// Each pipeline stage should own its own Sampler instance.
package sampler
