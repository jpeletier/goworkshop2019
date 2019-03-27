package streams_test

import (
	"bytes"
	"goworkshop/streams"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"testing"
)

func TestHexReader(t *testing.T) {
	sourceData := []byte("456E20756E206C75676172206465206C61204D616E6368612C206465206375796F206E6F6D627265206E6F2071756965726F2061636F726461726D652E2E2E")
	expectedResult := "En un lugar de la Mancha, de cuyo nombre no quiero acordarme..."

	var r io.Reader = bytes.NewReader(sourceData)
	r = streams.NewHexReader(r)

	var buf bytes.Buffer

	length, err := io.Copy(&buf, r)
	if err != nil {
		t.Fatal(err)
	}

	if length != int64(len(expectedResult)) {
		t.Fatalf("Expected length was %d, got %d", len(expectedResult), length)
	}

	if !bytes.Equal([]byte(expectedResult), buf.Bytes()) {
		t.Fatalf("Expected decoded result to be '%s', got '%s'", expectedResult, string(buf.Bytes()))
	}

}

func TestHexReaderInvalidChars(t *testing.T) {
	sourceData := []byte("ABCDEFGHIJKLMNOPQRST")

	var r io.Reader = bytes.NewReader(sourceData)
	r = streams.NewHexReader(r)

	_, err := ioutil.ReadAll(r)

	if err != streams.ErrInvalidHexChar {
		t.Fatal("Expected read operation to fail since source data contains invalid hex chars")
	}
}

func TestHexReaderOddNumberOfChars(t *testing.T) {
	sourceData := []byte("ABCD12345")

	var r io.Reader = bytes.NewReader(sourceData)
	r = streams.NewHexReader(r)

	_, err := ioutil.ReadAll(r)

	if err != streams.ErrMissingHexChar {
		t.Fatal("Expected read operation to fail since source data contains odd number of characters")
	}
}

func TestHexWriter(t *testing.T) {
	sourceData := "... no ha mucho tiempo que vivía un hidalgo de los de lanza en astillero, adarga antigua, rocín flaco y galgo corredor."
	expectedResult := "2E2E2E206E6F206861206D7563686F207469656D706F2071756520766976ED6120756E20686964616C676F206465206C6F73206465206C616E7A6120656E20617374696C6C65726F2C2061646172676120616E74696775612C20726F63ED6E20666C61636F20792067616C676F20636F727265646F722E"

	buf := new(bytes.Buffer)
	w := streams.NewHexWriter(buf)
	r := strings.NewReader(sourceData)

	_, err := io.Copy(w, r)
	if err != nil {
		t.Fatal(err)
	}

	result := buf.Bytes()

	if !bytes.Equal(result, []byte(expectedResult)) {
		t.Fatalf("Expected result to be '%s', got '%s'", expectedResult, string(result))
	}
}

func TestHexCombined(t *testing.T) {
	rndReader := rand.New(rand.NewSource(2019))
	expectedResult := new(bytes.Buffer)

	r := io.TeeReader(rndReader, expectedResult)

	buf := new(bytes.Buffer)

	hexWriter := streams.NewHexWriter(buf)

	_, err := io.CopyN(hexWriter, r, 10000)
	if err != nil {
		t.Fatal(err)
	}

	hexReader := streams.NewHexReader(buf)

	result, err := ioutil.ReadAll(hexReader)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, expectedResult.Bytes()) {
		log.Fatal("expected original array and result to be equal")
	}

}
