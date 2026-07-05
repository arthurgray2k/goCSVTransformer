package gocsvtransformer

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

var mapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]string)
	},
}

// CSVToJSON transforms a CSV stream from `r` to a JSON stream to `w`.
func (t *Transformer) CSVToJSON(ctx context.Context, r io.Reader, w io.Writer) error {
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

	if t.opts.Pretty {
		if _, err := w.Write([]byte("[\n")); err != nil {
			return err
		}
	} else {
		if _, err := w.Write([]byte("[")); err != nil {
			return err
		}
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	isFirst := true
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

		// Reuse map to minimize allocations
		obj := mapPool.Get().(map[string]string)

		for i, val := range record {
			key := fmt.Sprintf("Column%d", i+1)
			if t.opts.HeaderRow && i < len(headers) {
				key = headers[i]
			}
			obj[key] = val
		}

		// Separator
		if !isFirst {
			if t.opts.Pretty {
				if _, err := w.Write([]byte(",\n  ")); err != nil {
					return err
				}
			} else {
				if _, err := w.Write([]byte(",")); err != nil {
					return err
				}
			}
		} else {
			if t.opts.Pretty {
				if _, err := w.Write([]byte("  ")); err != nil {
					return err
				}
			}
			isFirst = false
		}

		if t.opts.Pretty {
			indent := strings.Repeat(" ", t.opts.Indent)
			if indent == "" {
				indent = "  "
			}
			b, err := json.MarshalIndent(obj, "  ", indent)
			if err != nil {
				return err
			}
			if _, err := w.Write(b); err != nil {
				return err
			}
		} else {
			b, err := json.Marshal(obj)
			if err != nil {
				return err
			}
			if _, err := w.Write(b); err != nil {
				return err
			}
		}

		// Clean up and return to pool
		for k := range obj {
			delete(obj, k)
		}
		mapPool.Put(obj)
	}

	// Write end of JSON array
	if t.opts.Pretty {
		if _, err := w.Write([]byte("\n]\n")); err != nil {
			return err
		}
	} else {
		if _, err := w.Write([]byte("]\n")); err != nil {
			return err
		}
	}

	return nil
}
