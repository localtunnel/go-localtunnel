package localtunnel

import (
	"fmt"
	"io"
	"io/ioutil"
)

// A simple limitedReader similar to io.LimitReader that also let's us know
// if we reached EOF
type limitedReader struct {
	reader    io.Reader
	maxBytes  int64
	lastError error
}

func (l *limitedReader) Read(p []byte) (int, error) {
	if l.maxBytes <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.maxBytes {
		p = p[0:l.maxBytes]
	}
	n, err := l.reader.Read(p)
	l.lastError = err
	l.maxBytes -= int64(n)
	return n, err
}

func (l *limitedReader) ReachedEOF() bool {
	return l.lastError == io.EOF
}

// Simple function that will read at-most maxSize and return an error if we
// didn't reach EOF.
func readAtmost(r io.ReadCloser, maxSize int64) ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	// Close the reader no matter what happens
	defer r.Close()

	// If r.maxSize is zero or less read the entire body regardless of length
	if maxSize <= 0 {
		return ioutil.ReadAll(r)
	}

	// Read at-most maxSize from body and check that we read it all
	reader := limitedReader{
		reader:   r,
		maxBytes: maxSize,
	}
	body, err := ioutil.ReadAll(&reader)
	if err != nil {
		return nil, err
	}
	if !reader.ReachedEOF() {
		return nil, fmt.Errorf("response larger than %d bytes", maxSize)
	}
	return body, nil
}
