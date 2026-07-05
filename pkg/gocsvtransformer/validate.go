package gocsvtransformer

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

var (
	ErrDuplicateHeader = errors.New("duplicate header found")
	ErrEmptyHeader     = errors.New("empty header found")
	ErrInvalidRowSize  = errors.New("invalid row size")
)

// ValidateCSV streams a CSV file and validates it.
// It checks for:
// - Valid CSV formatting (handled by encoding/csv)
// - Duplicate headers (if opts.HeaderRow is true)
// - Empty headers
// - Inconsistent column counts across rows
func (t *Transformer) ValidateCSV(ctx context.Context, r io.Reader) error {
	reader := csv.NewReader(r)
	reader.Comma = t.opts.Delimiter

	// We want to handle column count checks manually to provide our own descriptive errors if needed,
	// but csv.Reader natively checks FieldsPerRecord. We will let csv.Reader do its job.
	// FieldsPerRecord = 0 means it will use the first row's count for validation.
	reader.FieldsPerRecord = 0

	var headers []string
	var err error
	var lineCount int

	if t.opts.HeaderRow {
		headers, err = reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return errors.New("file is empty")
			}
			return fmt.Errorf("failed to read headers on line 1: %w", err)
		}
		lineCount++

		if err := validateHeaders(headers); err != nil {
			return err
		}
	}

	for {
		// Check context cancellation before reading next row
		if ctx.Err() != nil {
			return ctx.Err()
		}

		_, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // Done reading
			}

			// csv.ParseError already includes the line number and details
			var parseErr *csv.ParseError
			if errors.As(err, &parseErr) {
				if parseErr.Err == csv.ErrFieldCount {
					return fmt.Errorf("%w on line %d: %v", ErrInvalidRowSize, parseErr.Line, parseErr.Err)
				}
				return fmt.Errorf("csv parse error on line %d: %w", parseErr.Line, parseErr.Err)
			}
			return fmt.Errorf("failed to read csv on line %d: %w", lineCount+1, err)
		}
		lineCount++
	}

	return nil
}

func validateHeaders(headers []string) error {
	seen := make(map[string]bool)
	for i, h := range headers {
		if h == "" {
			return fmt.Errorf("%w at column %d", ErrEmptyHeader, i+1)
		}
		if seen[h] {
			return fmt.Errorf("%w: '%s' at column %d", ErrDuplicateHeader, h, i+1)
		}
		seen[h] = true
	}
	return nil
}
