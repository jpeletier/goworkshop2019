package streams

import (
	"fmt"
	"io"
)

type SpyWriter struct {
	innerWriter io.Writer
}

func NewSpyWriter(w io.Writer) *SpyWriter {
	return &SpyWriter{
		innerWriter: w,
	}
}

func (w *SpyWriter) Write(p []byte) (int, error) {
	i, err := w.innerWriter.Write(p)
	fmt.Printf("\nCaptured data: %s", string(p))
	return i, err
}
