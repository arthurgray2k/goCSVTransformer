package gocsvtransformer

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// NDJSONToCSV transforms an NDJSON stream from `r` to a CSV stream to `w`.
func (t *Transformer) NDJSONToCSV(ctx context.Context, r io.Reader, w io.Writer) error {
	decoder := json.NewDecoder(r)

	writer := csv.NewWriter(w)
	writer.Comma = t.opts.Delimiter
	defer writer.Flush()

	var headers []string
	isFirst := true
	lineCount := 0

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		var obj map[string]interface{}
		err := decoder.Decode(&obj)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // end of stream
			}
			return fmt.Errorf("failed to decode object on line %d: %w", lineCount+1, err)
		}
		lineCount++

		// Initialize headers on the first object
		if isFirst {
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
				case float64, bool:
					row[i] = fmt.Sprintf("%v", v)
				default:
					b, _ := json.Marshal(v)
					row[i] = string(b)
				}
			} else {
				row[i] = ""
			}
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row on line %d: %w", lineCount, err)
		}
	}

	return nil
}
