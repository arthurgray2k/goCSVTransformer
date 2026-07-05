package gocsvtransformer_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/arthurgray2k/gocsvtransformer/pkg/gocsvtransformer"
)

func TestTransformer_CSVToNDJSON(t *testing.T) {
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
				lines := strings.Split(strings.TrimSpace(out), "\n")
				if len(lines) != 2 {
					t.Fatalf("Expected 2 lines, got %d", len(lines))
				}
				var obj map[string]string
				if err := json.Unmarshal([]byte(lines[0]), &obj); err != nil {
					t.Fatalf("Failed to parse NDJSON line 1: %v", err)
				}
				if obj["id"] != "1" || obj["name"] != "alice" {
					t.Errorf("Unexpected row 1: %v", obj)
				}
			},
		},
		{
			name:    "without headers",
			input:   "1,alice\n2,bob\n",
			options: gocsvtransformer.Options{HeaderRow: false, Delimiter: ','},
			wantErr: false,
			verify: func(t *testing.T, out string) {
				lines := strings.Split(strings.TrimSpace(out), "\n")
				if len(lines) != 2 {
					t.Fatalf("Expected 2 lines, got %d", len(lines))
				}
				var obj map[string]string
				if err := json.Unmarshal([]byte(lines[0]), &obj); err != nil {
					t.Fatalf("Failed to parse NDJSON line 1: %v", err)
				}
				if obj["Column1"] != "1" || obj["Column2"] != "alice" {
					t.Errorf("Unexpected row 1: %v", obj)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := gocsvtransformer.New(tt.options)
			var buf bytes.Buffer
			err := transformer.CSVToNDJSON(context.Background(), strings.NewReader(tt.input), &buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("CSVToNDJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.verify != nil {
				tt.verify(t, buf.String())
			}
		})
	}
}
