package reader

import (
	"bufio"
	"bytes"
	"io"
)

// LineReader is an interface to read lines with optional Metadata.
type LineReader interface {
	ReadLine() (string, *Metadata, error)
}

// Metadata contains metadata that belongs to a line.
type Metadata struct {
	Package string
}

// LimitedLineReader reads lines from an io.Reader object with a configurable
// line size limit. Lines exceeding the limit will be truncated, but read
// completely from the underlying io.Reader.
type LimitedLineReader struct {
	r     *bufio.Reader
	limit int
}

var _ LineReader = &LimitedLineReader{}

// NewLimitedLineReader returns a LimitedLineReader to read lines from r with a
// maximum line size of limit.
func NewLimitedLineReader(r io.Reader, limit int) *LimitedLineReader {
	return &LimitedLineReader{r: bufio.NewReader(r), limit: limit}
}

// ReadLine returns the next line from the underlying reader. The length of the
// line will not exceed the configured limit. ReadLine either returns a line or
// it returns an error, never both.
func (r *LimitedLineReader) ReadLine() (string, *Metadata, error) {
	line, isPrefix, err := r.r.ReadLine()
	if err != nil {
		return "", nil, err
	}

	if !isPrefix {
		return string(line), nil, nil
	}

	// Line is incomplete, keep reading until we reach the end of the line.
	var buf bytes.Buffer
	buf.Write(line) // ignore err, always nil
	for isPrefix {
		line, isPrefix, err = r.r.ReadLine()
		if err != nil {
			return "", nil, err
		}

		if buf.Len() >= r.limit {
			// Stop writing to buf if we exceed the limit. We continue reading
			// however to make sure we consume the entire line.
			continue
		}

		buf.Write(line) // ignore err, always nil
	}

	if buf.Len() > r.limit {
		buf.Truncate(r.limit)
	}
	return buf.String(), nil, nil
}
