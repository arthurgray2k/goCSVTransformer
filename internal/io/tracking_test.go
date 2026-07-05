package io_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	internalio "github.com/arthurgray2k/gocsvtransformer/internal/io"
)

func TestTrackingReader(t *testing.T) {
	input := "hello world"
	r := strings.NewReader(input)
	tr := internalio.NewTrackingReader(r)

	buf := make([]byte, 5)
	n, err := tr.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected to read 5 bytes, got %d", n)
	}
	if tr.BytesRead() != 5 {
		t.Errorf("expected 5 bytes tracked, got %d", tr.BytesRead())
	}

	_, err = io.ReadAll(tr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tr.BytesRead() != int64(len(input)) {
		t.Errorf("expected %d bytes total, got %d", len(input), tr.BytesRead())
	}

	if tr.Elapsed() < 0 {
		t.Errorf("expected elapsed time to be positive, got %v", tr.Elapsed())
	}
}

func TestTrackingWriter(t *testing.T) {
	var buf bytes.Buffer
	tw := internalio.NewTrackingWriter(&buf)

	data := []byte("hello world")
	n, err := tw.Write(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(data) {
		t.Errorf("expected to write %d bytes, got %d", len(data), n)
	}
	if tw.BytesWritten() != int64(len(data)) {
		t.Errorf("expected %d bytes tracked, got %d", len(data), tw.BytesWritten())
	}

	// Sleep slightly to ensure elapsed time > 0 if timer resolution is coarse
	time.Sleep(1 * time.Millisecond)

	if tw.Elapsed() <= 0 {
		t.Errorf("expected elapsed time to be strictly positive, got %v", tw.Elapsed())
	}
}
