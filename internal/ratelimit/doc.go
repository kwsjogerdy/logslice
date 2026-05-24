// Package ratelimit provides a token-bucket rate limiter that controls how
// many log lines logslice emits per second.
//
// Usage:
//
//	limiter := ratelimit.New(1000) // allow at most 1 000 lines/s
//	for _, line := range lines {
//		if limiter.Allow() {
//			fmt.Println(line)
//		}
//	}
//
// A rate of 0 (or negative) disables limiting so that every line is allowed
// through, which is the default behaviour when the flag is not supplied.
package ratelimit
