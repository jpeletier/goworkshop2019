package streams

import (
	"errors"
	"io"
)

var ErrInvalidHexChar = errors.New("Invalid character in hex stream")
var ErrMissingHexChar = errors.New("Hex stream contains an odd number of characters")

type HexReader struct {
}

func NewHexReader(r io.Reader) io.Reader {
	return nil
}

func (r *HexReader) Read(p []byte) (int, error) {
	return 0, nil
}

type HexWriter struct {
}

func NewHexWriter(w io.Writer) io.Writer {
	return nil
}

func (w *HexWriter) Write(p []byte) (int, error) {
	return 0, nil
}
