package gocsvtransformer_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/arthurgray2k/gocsvtransformer/pkg/gocsvtransformer"
)

func TestTransformer_CSVToJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		options gocsvtransformer.Options
		wantErr bool
		verify  func(*testing.T, string)
	}{
		{
			name:    "with headers",
			input:   "id,name\n1,alice\n2,bob\n",
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ','},
			wantErr: false,
			verify: func(t *testing.T, out string) {
				var data []map[string]string
				if err := json.Unmarshal([]byte(out), &data); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				if len(data) != 2 {
					t.Fatalf("Expected 2 rows, got %d", len(data))
				}
				if data[0]["id"] != "1" || data[0]["name"] != "alice" {
					t.Errorf("Unexpected row 1: %v", data[0])
				}
			},
		},
		{
			name:    "without headers",
			input:   "1,alice\n2,bob\n",
			options: gocsvtransformer.Options{HeaderRow: false, Delimiter: ','},
			wantErr: false,
			verify: func(t *testing.T, out string) {
				var data []map[string]string
				if err := json.Unmarshal([]byte(out), &data); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				if len(data) != 2 {
					t.Fatalf("Expected 2 rows, got %d", len(data))
				}
				if data[0]["Column1"] != "1" || data[0]["Column2"] != "alice" {
					t.Errorf("Unexpected row 1: %v", data[0])
				}
			},
		},
		{
			name:    "pretty output",
			input:   "id,name\n1,alice\n",
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ',', Pretty: true, Indent: 2},
			wantErr: false,
			verify: func(t *testing.T, out string) {
				if !strings.Contains(out, "[\n  {\n") {
					t.Errorf("JSON does not look pretty: %s", out)
				}
				var data []map[string]string
				if err := json.Unmarshal([]byte(out), &data); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := gocsvtransformer.New(tt.options)
			var buf bytes.Buffer
			err := transformer.CSVToJSON(context.Background(), strings.NewReader(tt.input), &buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("CSVToJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.verify != nil {
				tt.verify(t, buf.String())
			}
		})
	}
}
