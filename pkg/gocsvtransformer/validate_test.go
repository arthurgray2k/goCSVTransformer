package gocsvtransformer_test

import (
	"context"
	"strings"
	"testing"

	"github.com/arthurgray2k/gocsvtransformer/pkg/gocsvtransformer"
)

func TestTransformer_ValidateCSV(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		options gocsvtransformer.Options
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid csv with headers",
			input:   "id,name,age\n1,alice,30\n2,bob,40\n",
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ','},
			wantErr: false,
		},
		{
			name:    "duplicate headers",
			input:   "id,name,name\n1,alice,30\n",
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ','},
			wantErr: true,
			errMsg:  "duplicate header",
		},
		{
			name:    "empty header",
			input:   "id,,age\n1,alice,30\n",
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ','},
			wantErr: true,
			errMsg:  "empty header",
		},
		{
			name:    "inconsistent column counts",
			input:   "id,name,age\n1,alice,30\n2,bob\n",
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: ','},
			wantErr: true,
			errMsg:  "invalid row size",
		},
		{
			name:    "custom delimiter",
			input:   "id|name|age\n1|alice|30\n",
			options: gocsvtransformer.Options{HeaderRow: true, Delimiter: '|'},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := gocsvtransformer.New(tt.options)
			err := transformer.ValidateCSV(context.Background(), strings.NewReader(tt.input))

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateCSV() error message = %v, want to contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}
