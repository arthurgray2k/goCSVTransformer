package gocsvtransformer

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// CSVToNDJSON transforms a CSV stream from `r` to an NDJSON stream to `w`.
func (t *Transformer) CSVToNDJSON(ctx context.Context, r io.Reader, w io.Writer) error {
	reader := csv.NewReader(r)
	reader.Comma = t.opts.Delimiter

	var headers []string
	var err error

	if t.opts.HeaderRow {
		headers, err = reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("failed to read header row: %w", err)
		}
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	lineCount := 0
	if t.opts.HeaderRow {
		lineCount = 1
	}

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("failed to read row %d: %w", lineCount+1, err)
		}
		lineCount++

		obj := mapPool.Get().(map[string]string)

		for i, val := range record {
			key := fmt.Sprintf("Column%d", i+1)
			if t.opts.HeaderRow && i < len(headers) {
				key = headers[i]
			}
			obj[key] = val
		}

		if err := encoder.Encode(obj); err != nil {
			return fmt.Errorf("failed to encode object on line %d: %w", lineCount, err)
		}

		for k := range obj {
			delete(obj, k)
		}
		mapPool.Put(obj)
	}

	return nil
}
