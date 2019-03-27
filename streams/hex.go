package streams

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
)

var ErrInvalidHexChar = errors.New("Invalid character in hex stream")
var ErrMissingHexChar = errors.New("Hex stream contains an odd number of characters")

type HexReader struct {
	ir io.Reader
}

func NewHexReader(r io.Reader) io.Reader {
	return &HexReader{
		ir: r,
	}
}

func (r *HexReader) Read(p []byte) (int, error) {
	d := make([]byte, len(p)*2)
	i, err := r.ir.Read(d)
	if err != nil {
		return i, err
	}

	if i%2 != 0 {
		return i, ErrMissingHexChar
	}

	d = d[:i]
	n, err := hex.Decode(p, d)
	if err != nil {
		err = ErrInvalidHexChar
	}
	return n, err
}

type HexWriter struct {
	iw io.Writer
}

func NewHexWriter(w io.Writer) io.Writer {
	return &HexWriter{
		iw: w,
	}
}

func (w *HexWriter) Write(p []byte) (int, error) {
	d := make([]byte, len(p)*2)
	hex.Encode(d, p)
	_, err := w.iw.Write(bytes.ToUpper(d))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
