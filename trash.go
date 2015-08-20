package trash

import "io"

func init() {
	Reader = &reader{}
}

// Reader provides a steady stream of trash (non-random bytes) when read from
var Reader io.Reader

type reader struct{}

func (r *reader) Read(p []byte) (int, error) {
	c := 0
	for i := 0; i < len(p); i++ {
		c++
		p[i] = 0xca
	}
	return c, nil
}
