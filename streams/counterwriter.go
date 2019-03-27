package streams

import (
	"io"
)

type CounterWriter struct {
	innerWriter io.Writer
	Length      int
}

func NewCounterWriter(w io.Writer) *CounterWriter {
	return &CounterWriter{
		innerWriter: w,
	}
}

func (w *CounterWriter) Write(p []byte) (int, error) {
	i, err := w.innerWriter.Write(p)
	w.Length += i
	return i, err
}
