package streams_test

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"goworkshop/streams"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestDumpFile(t *testing.T) {
	file, err := os.Open("sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownloadFile(t *testing.T) {
	r, err := http.Get("https://www.gutenberg.org/cache/epub/2000/pg2000.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	file, err := os.Create("download.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	length, err := io.Copy(file, r.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%d bytes written\n", length)

}

func TestCalcHash(t *testing.T) {
	file, err := os.Open("sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	hasher := sha256.New()

	_, err = io.Copy(hasher, file)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("The SHA256 hash of the file is %s\n", hex.EncodeToString(hasher.Sum(nil)))
}

func TestZipFile(t *testing.T) {
	srcFile, err := os.Open("sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create("sample.txt.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer dstFile.Close()

	zipper := gzip.NewWriter(dstFile)
	defer zipper.Close()

	_, err = io.Copy(zipper, srcFile)
	if err != nil {
		t.Fatal(err)
	}

}

func TestZipFileByteCount(t *testing.T) {
	srcFile, err := os.Open("sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create("sample.txt.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer dstFile.Close()

	counterWriter := streams.NewCounterWriter(dstFile)

	zipper := gzip.NewWriter(counterWriter)

	length, err := io.Copy(zipper, srcFile)
	if err != nil {
		t.Fatal(err)
	}

	zipper.Close()

	fmt.Printf("\n\nOriginal size: %d, compressed size: %d\n", length, counterWriter.Length)

}

func TestUnzipFile(t *testing.T) {
	srcFile, err := os.Open("sample.txt.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer srcFile.Close()

	unzipper, err := gzip.NewReader(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(os.Stdout, unzipper)
	if err != nil {
		t.Fatal(err)
	}

}

func TestDownloadAndHash(t *testing.T) {

	r, err := http.Get("https://www.gutenberg.org/cache/epub/2000/pg2000.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	file, err := os.Create("download.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	hasher := sha256.New()
	teeReader := io.TeeReader(r.Body, hasher)

	_, err = io.Copy(file, teeReader)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("The hash is %s\n", hex.EncodeToString(hasher.Sum(nil)))

}
