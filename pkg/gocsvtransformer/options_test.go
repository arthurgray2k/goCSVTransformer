package gocsvtransformer_test

import (
	"testing"

	"github.com/arthurgray2k/gocsvtransformer/pkg/gocsvtransformer"
)

func TestDefaultOptions(t *testing.T) {
	opts := gocsvtransformer.DefaultOptions()
	if opts.Delimiter != ',' {
		t.Errorf("expected delimiter ',', got %v", opts.Delimiter)
	}
	if !opts.HeaderRow {
		t.Errorf("expected HeaderRow to be true")
	}
	if opts.Pretty {
		t.Errorf("expected Pretty to be false")
	}
	if opts.Indent != 2 {
		t.Errorf("expected Indent to be 2, got %d", opts.Indent)
	}
}

func TestNew_EmptyDelimiter(t *testing.T) {
	opts := gocsvtransformer.Options{}
	tr := gocsvtransformer.New(opts)
	if tr == nil {
		t.Fatalf("expected transformer, got nil")
	}
}
