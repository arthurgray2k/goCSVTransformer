package gocsvtransformer_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/arthurgray2k/gocsvtransformer/pkg/gocsvtransformer"
)

func TestTransformer_NDJSONToCSV(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		options gocsvtransformer.Options
		wantErr bool
		verify  func(*testing.T, string)
	}{
		{
			name:    "basic objects",
			input:   "{\"id\": 1, \"name\": \"alice\"}\n{\"id\": 2, \"name\": \"bob\"}\n",
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ','},
			wantErr: false,
			verify: func(t *testing.T, out string) {
				lines := strings.Split(strings.TrimSpace(out), "\n")
				if len(lines) != 3 { // header + 2 records
					t.Fatalf("Expected 3 lines (header + 2 rows), got %d", len(lines))
				}
				if !strings.Contains(out, "1") || !strings.Contains(out, "alice") {
					t.Errorf("Output missing data: %s", out)
				}
			},
		},
		{
			name:    "missing keys in subsequent rows",
			input:   "{\"id\": 1, \"name\": \"alice\"}\n{\"id\": 2}\n",
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ','},
			wantErr: false,
			verify: func(t *testing.T, out string) {
				lines := strings.Split(strings.TrimSpace(out), "\n")
				if len(lines) != 3 {
					t.Fatalf("Expected 3 lines, got %d", len(lines))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := gocsvtransformer.New(tt.options)
			var buf bytes.Buffer
			err := transformer.NDJSONToCSV(context.Background(), strings.NewReader(tt.input), &buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("NDJSONToCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.verify != nil {
				tt.verify(t, buf.String())
			}
		})
	}
}
