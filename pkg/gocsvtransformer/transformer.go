package gocsvtransformer

import (
	"errors"
)

// ErrNotImplemented is returned for methods that are planned but not yet implemented.
var ErrNotImplemented = errors.New("not implemented")

// Transformer represents the main execution engine for CSV and JSON data transformations.
type Transformer struct {
	opts Options
}

// New creates a new Transformer with the provided Options.
func New(opts Options) *Transformer {
	// Fallback to comma if empty
	if opts.Delimiter == 0 {
		opts.Delimiter = ','
	}
	return &Transformer{
		opts: opts,
	}
}
