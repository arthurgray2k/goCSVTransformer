package io

import (
	"io"
	"time"
)

// TrackingReader wraps an io.Reader and tracks the number of bytes read.
type TrackingReader struct {
	r     io.Reader
	bytes int64
	start time.Time
}

// NewTrackingReader creates a new TrackingReader.
func NewTrackingReader(r io.Reader) *TrackingReader {
	return &TrackingReader{
		r:     r,
		start: time.Now(),
	}
}

// Read implements io.Reader.
func (t *TrackingReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	t.bytes += int64(n)
	return n, err
}

// BytesRead returns the total number of bytes read.
func (t *TrackingReader) BytesRead() int64 {
	return t.bytes
}

// Elapsed returns the time elapsed since the reader was created.
func (t *TrackingReader) Elapsed() time.Duration {
	return time.Since(t.start)
}

// TrackingWriter wraps an io.Writer and tracks the number of bytes written.
type TrackingWriter struct {
	w     io.Writer
	bytes int64
	start time.Time
}

// NewTrackingWriter creates a new TrackingWriter.
func NewTrackingWriter(w io.Writer) *TrackingWriter {
	return &TrackingWriter{
		w:     w,
		start: time.Now(),
	}
}

// Write implements io.Writer.
func (t *TrackingWriter) Write(p []byte) (n int, err error) {
	n, err = t.w.Write(p)
	t.bytes += int64(n)
	return n, err
}

// BytesWritten returns the total number of bytes written.
func (t *TrackingWriter) BytesWritten() int64 {
	return t.bytes
}

// Elapsed returns the time elapsed since the writer was created.
func (t *TrackingWriter) Elapsed() time.Duration {
	return time.Since(t.start)
}
