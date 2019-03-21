package swarmcas

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/swarm/storage"
)

type swarmClientMock struct {
	table map[string][]byte
}

func NewSwarmClientMock() SwarmClientRaw {
	return &swarmClientMock{
		table: make(map[string][]byte),
	}
}

func (sm *swarmClientMock) UploadRaw(r io.Reader, size int64, encrypt bool) (string, error) {
	if encrypt {
		panic("mock does not support encryption")
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	key, err := SwarmHash(data)
	if err != nil {
		return "", err
	}

	sm.table[key] = data
	return key, nil
}

func (sm *swarmClientMock) DownloadRaw(key string) (io.ReadCloser, bool, error) {
	data := sm.table[key]
	if data == nil {
		return nil, false, errors.New("cannot find key")
	}
	return ioutil.NopCloser(bytes.NewReader(data)), false, nil
}

func SwarmHash(data []byte) (string, error) {
	fileStore := storage.NewFileStore(&storage.FakeChunkStore{}, storage.NewFileStoreParams())
	key, _, err := fileStore.Store(context.TODO(), bytes.NewReader(data), int64(len(data)), false)
	if err != nil {
		return "", err
	}

	return key.Hex(), nil
}
