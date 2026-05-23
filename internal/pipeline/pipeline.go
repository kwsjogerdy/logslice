// Package pipeline wires together the reader, parser, filter, formatter,
// and writer into a single processing pipeline for log lines.
package pipeline

import (
	"context"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/reader"
)

// Config holds the configuration for a pipeline run.
type Config struct {
	Include    []string
	Exclude    []string
	MinLevel   string
	Colorize   bool
	InputFile  string
	OutputFile string
}

// Pipeline processes log lines from a reader through filter/format to a writer.
type Pipeline struct {
	reader    *reader.Reader
	filter    *filter.Filter
	formatter *formatter.Formatter
	writer    *output.Writer
}

// New constructs a Pipeline from the provided Config.
func New(cfg Config) (*Pipeline, error) {
	r, err := openReader(cfg.InputFile)
	if err != nil {
		return nil, err
	}

	w, err := openWriter(cfg.OutputFile)
	if err != nil {
		return nil, err
	}

	minLevel := formatter.ParseLevel(cfg.MinLevel)

	return &Pipeline{
		reader:    r,
		filter:    filter.New(cfg.Include, cfg.Exclude),
		formatter: formatter.New(minLevel, cfg.Colorize),
		writer:    w,
	}, nil
}

// Run reads all lines, applies the filter and formatter, and writes results.
// It respects context cancellation.
func (p *Pipeline) Run(ctx context.Context) (Stats, error) {
	var stats Stats
	for {
		select {
		case <-ctx.Done():
			return stats, ctx.Err()
		default:
		}

		line, err := p.reader.ReadLine()
		if err != nil {
			break
		}

		stats.Read++
		entry := parser.Parse(line)

		if !p.filter.Match(entry) {
			stats.Filtered++
			continue
		}

		formatted, ok := p.formatter.Format(entry)
		if !ok {
			stats.Filtered++
			continue
		}

		if writeErr := p.writer.WriteLine(formatted); writeErr != nil {
			return stats, writeErr
		}
		stats.Written++
	}
	return stats, nil
}

func openReader(path string) (*reader.Reader, error) {
	if path == "" {
		return reader.Stdin(), nil
	}
	return reader.OpenFile(path)
}

func openWriter(path string) (*output.Writer, error) {
	if path == "" {
		return output.Stdout(), nil
	}
	return output.OpenFile(path)
}
