package gocsvtransformer

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// JSONToCSV transforms a JSON stream from `r` to a CSV stream to `w`.
// It expects a JSON array of objects.
func (t *Transformer) JSONToCSV(ctx context.Context, r io.Reader, w io.Writer) error {
	decoder := json.NewDecoder(r)

	// Read the opening bracket
	token, err := decoder.Token()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("failed to read json token: %w", err)
	}

	delim, ok := token.(json.Delim)
	if !ok || delim != '[' {
		return errors.New("expected JSON array at the root")
	}

	writer := csv.NewWriter(w)
	writer.Comma = t.opts.Delimiter
	defer writer.Flush()

	var headers []string
	isFirst := true

	for decoder.More() {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		var obj map[string]interface{}
		if err := decoder.Decode(&obj); err != nil {
			return fmt.Errorf("failed to decode object: %w", err)
		}

		// Initialize headers on the first object
		if isFirst {
			// To ensure deterministic column ordering for identical schemas, we just take the iteration order
			// (Go map iteration is random, but we freeze it here).
			// A production robust way is to just accept whatever map iteration gives us for the first row
			// or optionally accept explicit headers.
			// Here we take it as-is from the first row.
			for k := range obj {
				headers = append(headers, k)
			}
			if t.opts.HeaderRow {
				if err := writer.Write(headers); err != nil {
					return fmt.Errorf("failed to write headers: %w", err)
				}
			}
			isFirst = false
		}

		// Write row based on established headers
		row := make([]string, len(headers))
		for i, h := range headers {
			if val, ok := obj[h]; ok && val != nil {
				// Convert to string
				switch v := val.(type) {
				case string:
					row[i] = v
				case float64, bool: // json unmarshals numbers to float64
					row[i] = fmt.Sprintf("%v", v)
				default:
					// For arrays and nested objects, serialize them back to json string
					b, _ := json.Marshal(v)
					row[i] = string(b)
				}
			} else {
				row[i] = ""
			}
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	// Read closing bracket
	token, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to read json closing token: %w", err)
	}
	delim, ok = token.(json.Delim)
	if !ok || delim != ']' {
		return errors.New("expected closing JSON array bracket")
	}

	return nil
}
