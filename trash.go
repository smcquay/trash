// Package trash provides simple readers with meaningless output.
package trash

import (
	"io"
	"time"
)

func init() {
	Reader = &reader{}
}

// Reader provides a steady stream of trash (non-random bytes) when read from
var Reader io.Reader

// TimeoutReader returns a reader that returns io.EOF after dur.
func TimeoutReader(dur time.Duration) io.Reader {
	return &timeoutReader{timeout: time.Now().Add(dur)}
}

type reader struct{}

func (r *reader) Read(p []byte) (int, error) {
	c := 0
	var err error
	for i := 0; i < len(p); i++ {
		c++
		p[i] = 0xca
	}
	return c, err
}

type timeoutReader struct {
	timeout time.Time
}

func (tor *timeoutReader) Read(p []byte) (int, error) {
	c := 0
	var err error
	for i := 0; i < len(p); i++ {
		c++
		p[i] = 0xca
	}
	if time.Now().After(tor.timeout) {
		err = io.EOF
	}
	return c, err
}
