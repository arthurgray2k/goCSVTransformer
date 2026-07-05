package gocsvtransformer_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/arthurgray2k/gocsvtransformer/pkg/gocsvtransformer"
)

func TestTransformer_JSONToCSV(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		options gocsvtransformer.Options
		wantErr bool
		verify  func(*testing.T, string)
	}{
		{
			name:    "basic array of objects",
			input:   `[{"id": 1, "name": "alice"}, {"id": 2, "name": "bob"}]`,
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ','},
			wantErr: false,
			verify: func(t *testing.T, out string) {
				lines := strings.Split(strings.TrimSpace(out), "\n")
				if len(lines) != 3 {
					t.Fatalf("Expected 3 lines (header + 2 rows), got %d", len(lines))
				}
				// Note: map iteration order is random, so we just check if it contains the values.
				if !strings.Contains(out, "1") || !strings.Contains(out, "alice") {
					t.Errorf("Output missing data: %s", out)
				}
			},
		},
		{
			name:    "missing keys in subsequent rows",
			input:   `[{"id": 1, "name": "alice"}, {"id": 2}]`,
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ','},
			wantErr: false,
			verify: func(t *testing.T, out string) {
				lines := strings.Split(strings.TrimSpace(out), "\n")
				if len(lines) != 3 {
					t.Fatalf("Expected 3 lines, got %d", len(lines))
				}
			},
		},
		{
			name:    "nested object",
			input:   `[{"id": 1, "details": {"age": 30}}]`,
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ','},
			wantErr: false,
			verify: func(t *testing.T, out string) {
				if !strings.Contains(out, `"{""age"":30}"`) { // csv escapes double quotes
					t.Errorf("Output missing serialized nested object: %s", out)
				}
			},
		},
		{
			name:    "not an array",
			input:   `{"id": 1}`,
			options: gocsvtransformer.Options{HeaderRow: true},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := gocsvtransformer.New(tt.options)
			var buf bytes.Buffer
			err := transformer.JSONToCSV(context.Background(), strings.NewReader(tt.input), &buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("JSONToCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.verify != nil {
				tt.verify(t, buf.String())
			}
		})
	}
}
