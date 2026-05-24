// Package config defines the Config struct that aggregates all runtime
// settings for logslice and provides helpers for populating it from
// command-line flags.
//
// Typical usage:
//
//	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
//	cfg, err := config.FromFlags(fs, os.Args[1:])
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The zero value of Config is not useful; use config.Default() to obtain
// a Config with sensible defaults before customising individual fields.
package config
